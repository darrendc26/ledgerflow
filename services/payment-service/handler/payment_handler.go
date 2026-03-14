package handler

import (
	"context"
	"fmt"

	"ledgerflow/proto/ledgerpb"
	pb "ledgerflow/proto/paymentpb"

	payment_service "ledgerflow/services/payment-service/payment_service"
)

type PaymentHandler struct {
	pb.UnimplementedPaymentServiceServer
	paymentService *payment_service.PaymentService
	ledgerClient   ledgerpb.LedgerServiceClient
}

func NewPaymentHandler(
	paymentService *payment_service.PaymentService,
	ledgerClient ledgerpb.LedgerServiceClient,
) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		ledgerClient:   ledgerClient,
	}
}

func (h *PaymentHandler) CreatePayment(
	ctx context.Context,
	req *pb.CreatePaymentRequest,
) (*pb.CreatePaymentResponse, error) {
	ledgerResp, err := h.ledgerClient.Transfer(ctx, &ledgerpb.TransferRequest{
		SenderAccount:   req.SenderAccount,
		ReceiverAccount: req.ReceiverAccount,
		Amount:          req.Amount,
		ReferenceId:     "",
	})

	if err != nil {
		return nil, err
	}
	fmt.Println(ledgerResp.Status)

	payment, err := h.paymentService.CreatePayment(
		req.SenderAccount,
		req.ReceiverAccount,
		req.Amount,
		req.Currency,
	)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePaymentResponse{
		PaymentId: payment.ID,
		Status:    payment.Status,
	}, nil
}
