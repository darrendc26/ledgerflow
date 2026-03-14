package handler

import (
	"context"

	pb "ledgerflow/proto/ledgerpb"
	"ledgerflow/services/ledger/service"
)

type LedgerHandler struct {
	pb.UnimplementedLedgerServiceServer
	service *service.LedgerService
}

func NewLedgerHandler(service *service.LedgerService) *LedgerHandler {
	return &LedgerHandler{service: service}
}

func (h *LedgerHandler) Transfer(ctx context.Context, req *pb.TransferRequest) (*pb.TransferResponse, error) {
	transfer, err := h.service.CreateTransfer(ctx, req.SenderAccount, req.ReceiverAccount, req.Amount, req.ReferenceId)
	if err != nil {
		return nil, err
	}
	return &pb.TransferResponse{
		TransactionId: transfer.TransactionID,
		Status:        transfer.Status,
	}, nil
}
