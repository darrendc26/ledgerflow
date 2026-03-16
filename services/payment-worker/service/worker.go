package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"ledgerflow/infra/prometheus"
	ledgerpb "ledgerflow/proto/ledgerpb"
	"ledgerflow/services/payment-service/kafka"
	service "ledgerflow/services/payment-service/payment_service"

	kafkago "github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
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
		producer:       producer,
	}
}

func (w *Worker) Start() {

	metrics := prometheus.NewPrometheus()
	metrics.Start(":2113")

	tracer := otel.Tracer("payment-worker")

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

		start := time.Now()

		// create trace span per payment
		ctx, span := tracer.Start(context.Background(), "process-payment")

		// check payment status (idempotency)
		paymentStatus, err := w.paymentService.GetPendingPayment(event.PaymentID)
		if err != nil {
			log.Println("Payment not found:", err)
			span.End()
			continue
		}

		if paymentStatus != "created" {
			log.Println("Payment already processed:", event.PaymentID)
			span.End()
			continue
		}

		// retry ledger transfer
		for i := 0; i < 3; i++ {

			ledgerCtx, ledgerSpan := tracer.Start(ctx, "ledger-transfer")

			_, err = w.ledgerClient.Transfer(ledgerCtx, &ledgerpb.TransferRequest{
				SenderAccount:   event.SenderAccount,
				ReceiverAccount: event.ReceiverAccount,
				Amount:          event.Amount,
				ReferenceId:     event.PaymentID,
			})

			ledgerSpan.End()

			if err == nil {
				break
			}

			log.Println("Retry attempt", i+1, "failed:", err)

			time.Sleep(2 * time.Second)
		}

		// failure after retries
		if err != nil {

			log.Println("Payment failed:", err)

			err := w.producer.PublishDLQ(&kafka.PaymentEvent{
				PaymentID:       event.PaymentID,
				SenderAccount:   event.SenderAccount,
				ReceiverAccount: event.ReceiverAccount,
				Amount:          event.Amount,
			})

			if err != nil {
				log.Println("Failed to publish to DLQ:", err)
			}

			log.Println("Sent to DLQ")

			w.paymentService.UpdateStatus(event.PaymentID, "failed")

			metrics.PaymentsFailed.Inc()

			span.End()
			continue
		}

		// success
		err = w.paymentService.UpdateStatus(event.PaymentID, "completed")
		if err != nil {
			log.Println("Failed to update payment status:", err)
		}

		log.Println("Payment completed:", event.PaymentID)

		metrics.PaymentsProcessed.Inc()

		metrics.PaymentsLatency.Observe(time.Since(start).Seconds())

		span.End()
	}
}
