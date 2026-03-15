package main

import (
	"log"

	"ledgerflow/pkg/db"
	ledgerpb "ledgerflow/proto/ledgerpb"
	kafka "ledgerflow/services/payment-service/kafka"
	payment_service "ledgerflow/services/payment-service/payment_service"
	"ledgerflow/services/payment-service/repository"
	"ledgerflow/services/payment-worker/service"

	grpc "google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	ledgerClient := ledgerpb.NewLedgerServiceClient(conn)
	repo := repository.NewPaymentRepository(db.NewPostgresPool())
	paymentService := payment_service.NewPaymentService(repo)
	producer := kafka.NewProducer()
	worker := service.NewWorker(ledgerClient, paymentService, producer)
	log.Println("Payment worker started")

	worker.Start()
}
