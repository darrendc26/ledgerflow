package handler

import (
	"context"
	pb "ledgerflow/proto/paymentpb"
	"ledgerflow/services/payment-service/service"
)

type PaymentHandler struct {
	pb.UnimplementedPaymentServiceServer
	service *service.PaymentService
}

func NewPaymentHandler(service *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		service: service,
	}
}

func (h *PaymentHandler) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	payment, err := h.service.CreatePayment(req.UserId, req.Amount, req.Currency)
	if err != nil {
		return nil, err
	}
	return &pb.CreatePaymentResponse{
		PaymentId: payment.ID,
		Status:    payment.Status,
	}, nil
}
