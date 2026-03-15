package handler

import (
	"context"

	pb "ledgerflow/proto/paymentpb"

	producer "ledgerflow/services/payment-service/kafka"
	payment_service "ledgerflow/services/payment-service/payment_service"
)

type PaymentHandler struct {
	pb.UnimplementedPaymentServiceServer
	paymentService *payment_service.PaymentService
	producer       *producer.Producer
}

func NewPaymentHandler(
	paymentService *payment_service.PaymentService,
	producer *producer.Producer,
) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		producer:       producer,
	}
}

func (h *PaymentHandler) CreatePayment(
	ctx context.Context,
	req *pb.CreatePaymentRequest,
) (*pb.CreatePaymentResponse, error) {

	payment, err := h.paymentService.CreatePayment(
		req.SenderAccount,
		req.ReceiverAccount,
		req.Amount,
		req.Currency,
	)
	if err != nil {
		return nil, err
	}

	err = h.producer.Publish(&producer.PaymentEvent{
		PaymentID:       payment.ID,
		SenderAccount:   req.SenderAccount,
		ReceiverAccount: req.ReceiverAccount,
		Amount:          req.Amount,
	})

	if err != nil {
		h.paymentService.UpdateStatus(payment.ID, "failed")
		return nil, err
	}

	return &pb.CreatePaymentResponse{
		PaymentId: payment.ID,
		Status:    "processing",
	}, nil
}
