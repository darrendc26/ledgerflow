package service

import (
	"ledgerflow/services/payment-service/model"
	"ledgerflow/services/payment-service/repository"

	"github.com/google/uuid"
)

type PaymentService struct {
	repository *repository.PaymentRepository
}

func NewPaymentService(repository *repository.PaymentRepository) *PaymentService {
	return &PaymentService{
		repository: repository,
	}
}

func (s *PaymentService) CreatePayment(senderAccount string, receiverAccount string, amount int64, currency string) (*model.Payment, error) {
	payment := &model.Payment{
		ID:              uuid.New().String(),
		SenderAccount:   senderAccount,
		ReceiverAccount: receiverAccount,
		Amount:          amount,
		Currency:        currency,
		Status:          "created",
	}
	err := s.repository.CreatePayment(payment)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *PaymentService) UpdateStatus(paymentID string, status string) error {
	return s.repository.UpdateStatus(paymentID, status)
}

func (s *PaymentService) GetPendingPayment(paymentID string) (string, error) {
	paymentID, err := s.repository.GetPendingPayment(paymentID)
	if err != nil {
		return "", err
	}
	return paymentID, nil
}
