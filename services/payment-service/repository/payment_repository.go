package repository

import (
	"context"
	"ledgerflow/services/payment-service/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

func (r *PaymentRepository) CreatePayment(payment *model.Payment) error {

	query := `INSERT INTO payments (id, sender_account, receiver_account, amount, currency, status)
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(context.Background(),
		query,
		payment.ID,
		payment.SenderAccount,
		payment.ReceiverAccount,
		payment.Amount,
		payment.Currency,
		payment.Status,
	)

	return err
}

func (r *PaymentRepository) UpdateStatus(id string, status string) error {

	query := `UPDATE payments SET status = $1 WHERE id = $2`

	_, err := r.db.Exec(context.Background(),
		query,
		status,
		id,
	)

	return err
}

func (r *PaymentRepository) GetPendingPayment(paymentID string) (string, error) {

	query := `
		SELECT status
		FROM payments
		WHERE id=$1
	`

	err := r.db.QueryRow(
		context.Background(),
		query,
		paymentID,
	).Scan(&paymentID)

	if err != nil {
		return "", err
	}

	return paymentID, nil
}
