package main

import (
	"ledgerflow/pkg/db"
	"ledgerflow/services/ledger/handler"
	"ledgerflow/services/ledger/repository"
	"ledgerflow/services/ledger/server"
	"ledgerflow/services/ledger/service"
)

func main() {
	db.NewPostgresPool()
	repo := repository.NewLedgerRepository(db.NewPostgresPool())
	service := service.NewLedgerService(repo)
	handler := handler.NewLedgerHandler(service)
	server.StartGrpcServer(handler)
}
