package main

import (
	"ledgerflow/pkg/db"
	"ledgerflow/services/api-gateway/clients"
	"ledgerflow/services/payment-service/handler"
	payment_service "ledgerflow/services/payment-service/payment_service"
	"ledgerflow/services/payment-service/repository"
	"ledgerflow/services/payment-service/server"
)

func main() {
	db.NewPostgresPool()
	repo := repository.NewPaymentRepository(db.NewPostgresPool())
	service := payment_service.NewPaymentService(repo)

	ledger_client, _ := clients.NewLedgerClient()

	handler := handler.NewPaymentHandler(service, ledger_client)

	server.StartGrpcServer(handler)
}
