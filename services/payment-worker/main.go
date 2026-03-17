package main

import (
	"log"

	"context"
	"ledgerflow/pkg/db"
	"ledgerflow/pkg/telemetry"
	ledgerpb "ledgerflow/proto/ledgerpb"
	kafka "ledgerflow/services/payment-service/kafka"
	payment_service "ledgerflow/services/payment-service/payment_service"
	"ledgerflow/services/payment-service/repository"
	"ledgerflow/services/payment-worker/service"

	grpc "google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("ledger-service:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	tp, err := telemetry.InitTracer("payment-worker")
	if err != nil {
		log.Fatal(err)
	}
	defer tp.Shutdown(context.Background())

	ledgerClient := ledgerpb.NewLedgerServiceClient(conn)
	repo := repository.NewPaymentRepository(db.NewPostgresPool())
	paymentService := payment_service.NewPaymentService(repo)
	producer := kafka.NewProducer()
	worker := service.NewWorker(ledgerClient, paymentService, producer)
	log.Println("Payment worker started")

	worker.Start()
}
