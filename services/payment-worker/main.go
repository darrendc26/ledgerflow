package main

import (
	"log"

	"ledgerflow/pkg/db"
	ledgerpb "ledgerflow/proto/ledgerpb"
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
	worker := service.NewWorker(ledgerClient, paymentService)
	log.Println("Payment worker started")

	worker.Start()
}
