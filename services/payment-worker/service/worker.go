package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	ledgerpb "ledgerflow/proto/ledgerpb"
	"ledgerflow/services/payment-service/kafka"
	service "ledgerflow/services/payment-service/payment_service"

	kafkago "github.com/segmentio/kafka-go"
)

type Worker struct {
	reader         *kafkago.Reader
	ledgerClient   ledgerpb.LedgerServiceClient
	paymentService *service.PaymentService
	producer       *kafka.Producer
}

func NewWorker(
	ledgerClient ledgerpb.LedgerServiceClient,
	paymentService *service.PaymentService,
	producer *kafka.Producer,
) *Worker {

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

		}

		var event kafka.PaymentEvent

		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Println("JSON decode error:", err)
		}

		log.Println("Processing payment:", event.PaymentID)

		paymentStatus, err := w.paymentService.GetPendingPayment(event.PaymentID)
		if err != nil {
			log.Println("Payment not found:", err)
			continue
		}
		// log.Println("Payment status:", paymentStatus)
		if paymentStatus != "created" {
			log.Println("Payment already processed:", event.PaymentID)
			break
		}

		for i := 0; i < 3; i++ {

			_, err = w.ledgerClient.Transfer(context.Background(), &ledgerpb.TransferRequest{
				SenderAccount:   event.SenderAccount,
				ReceiverAccount: event.ReceiverAccount,
				Amount:          event.Amount,
				ReferenceId:     event.PaymentID,
			})

			if err == nil {
				break
			}

			log.Println("Retry attempt", i+1, "failed:", err)

			time.Sleep(2 * time.Second)
		}

		if err != nil {
			log.Println("Payment failed:", err)
			log.Println("Sending to DLQ")

			w.producer.PublishDLQ(&kafka.PaymentEvent{
				PaymentID:       event.PaymentID,
				SenderAccount:   event.SenderAccount,
				ReceiverAccount: event.ReceiverAccount,
				Amount:          event.Amount,
			})

			w.paymentService.UpdateStatus(event.PaymentID, "failed")

			continue
		}

		err = w.paymentService.UpdateStatus(event.PaymentID, "completed")
		log.Println("Payment completed:", event.PaymentID)
		if err != nil {
			log.Println("Failed to update payment status:", err)
		}
	}
}
