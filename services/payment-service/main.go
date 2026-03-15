package main

import (
	"ledgerflow/pkg/db"
	"ledgerflow/services/payment-service/handler"
	"ledgerflow/services/payment-service/kafka"
	payment_service "ledgerflow/services/payment-service/payment_service"
	"ledgerflow/services/payment-service/repository"
	"ledgerflow/services/payment-service/server"
)

func main() {
	db.NewPostgresPool()
	repo := repository.NewPaymentRepository(db.NewPostgresPool())
	service := payment_service.NewPaymentService(repo)

	// ledger_client, _ := clients.NewLedgerClient()

	producer := kafka.NewProducer()
	handler := handler.NewPaymentHandler(service, producer)

	server.StartGrpcServer(handler)
}
