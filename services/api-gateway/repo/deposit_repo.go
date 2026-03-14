package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DepositRepository struct {
	db *pgxpool.Pool
}

func NewDepositRepository(db *pgxpool.Pool) *DepositRepository {
	return &DepositRepository{db: db}
}

func (r *DepositRepository) CreateDeposit(
	ctx context.Context,
	depositAccount string,
	amount int64,
	currency string,
) (string, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	txID := uuid.New().String()

	_, err = tx.Exec(
		ctx,
		`INSERT INTO transactions (id) VALUES ($1)`,
		txID,
	)
	if err != nil {
		return "", err
	}

	_, err = tx.Exec(
		ctx,
		`INSERT INTO ledger_entries (transaction_id, account_id, amount, type)
		 VALUES ($1,$2,$3,$4)`,
		txID,
		depositAccount,
		amount,
		"credit",
	)
	if err != nil {
		return "", err
	}

	_, err = tx.Exec(
		ctx,
		`UPDATE accounts SET balance = balance + $1 WHERE id = $2`,
		amount,
		depositAccount,
	)
	if err != nil {
		return "", err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", err
	}

	return txID, nil
}
