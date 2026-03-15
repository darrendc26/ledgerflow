package service

import (
	"context"
	"encoding/json"
	"log"

	ledgerpb "ledgerflow/proto/ledgerpb"
	"ledgerflow/services/payment-service/kafka"
	service "ledgerflow/services/payment-service/payment_service"

	kafkago "github.com/segmentio/kafka-go"
)

type Worker struct {
	reader         *kafkago.Reader
	ledgerClient   ledgerpb.LedgerServiceClient
	paymentService *service.PaymentService
}

func NewWorker(ledgerClient ledgerpb.LedgerServiceClient, paymentService *service.PaymentService) *Worker {

	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "payments",
		GroupID: "payments-group",
	})

	return &Worker{
		reader:         reader,
		ledgerClient:   ledgerClient,
		paymentService: paymentService,
	}
}

func (w *Worker) Start() {

	for {

		msg, err := w.reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Kafka read error:", err)
			continue
		}

		var event kafka.PaymentEvent

		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Println("JSON decode error:", err)
			continue
		}

		log.Println("Processing payment:", event.PaymentID)

		paymentStatus, err := w.paymentService.GetPendingPayment(event.PaymentID)
		if err != nil {
			log.Println("Payment not found:", err)
			continue
		}
		log.Println("Payment status:", paymentStatus)
		if paymentStatus != "created" {
			log.Println("Payment already processed:", event.PaymentID)
			continue
		}

		_, err = w.ledgerClient.Transfer(context.Background(), &ledgerpb.TransferRequest{
			SenderAccount:   event.SenderAccount,
			ReceiverAccount: event.ReceiverAccount,
			Amount:          event.Amount,
			ReferenceId:     event.PaymentID,
		})

		if err != nil {
			log.Println("Ledger transfer failed:", err)
			w.paymentService.UpdateStatus(event.PaymentID, "failed")
			continue
		}
		w.paymentService.UpdateStatus(event.PaymentID, "completed")
		log.Println("Payment completed:", event.PaymentID)
	}
}
