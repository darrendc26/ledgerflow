package service

import (
	"context"
	"ledgerflow/services/ledger/model"
	"ledgerflow/services/ledger/repository"
)

type LedgerService struct {
	repo *repository.LedgerRepository
}

func NewLedgerService(repo *repository.LedgerRepository) *LedgerService {
	return &LedgerService{repo: repo}
}

func (s *LedgerService) CreateTransfer(ctx context.Context, senderAccount string, receiverAccount string, amount int64, referenceID string) (*model.Ledger, error) {
	ledger := &model.Ledger{
		SenderAccount:   senderAccount,
		ReceiverAccount: receiverAccount,
		Amount:          amount,
		ReferenceID:     referenceID,
	}
	err := s.repo.CreateTransfer(ledger)
	if err != nil {
		return nil, err
	}
	return ledger, nil
}
