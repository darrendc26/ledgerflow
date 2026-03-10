package main

import (
	"ledgerflow/pkg/db"
	"ledgerflow/services/payment-service/handler"
	"ledgerflow/services/payment-service/repository"
	"ledgerflow/services/payment-service/server"
	"ledgerflow/services/payment-service/service"
)

func main() {
	db.NewPostgresPool()
	repo := repository.NewPaymentRepository(db.NewPostgresPool())
	service := service.NewPaymentService(repo)
	handler := handler.NewPaymentHandler(service)
	server.StartGrpcServer(handler)
}
