package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"ledgerflow/services/ledger/model"
)

type LedgerRepository struct {
	db *pgxpool.Pool
}

func NewLedgerRepository(db *pgxpool.Pool) *LedgerRepository {
	return &LedgerRepository{db: db}
}

func (r *LedgerRepository) CreateTransfer(
	ledger *model.Ledger,
) error {

	ctx := context.Background()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	// Lock sender account row
	var balance int64
	err = tx.QueryRow(
		ctx,
		`SELECT balance FROM accounts WHERE id = $1 FOR UPDATE`,
		ledger.SenderAccount,
	).Scan(&balance)

	if err != nil {
		tx.Rollback(ctx)
		fmt.Println("Transfer error")
		return err
	}

	if balance < ledger.Amount {
		tx.Rollback(ctx)
		return fmt.Errorf("insufficient funds")
	}

	txID := uuid.New().String()

	// Create transaction record
	_, err = tx.Exec(
		ctx,
		`INSERT INTO transactions (id, reference_id) VALUES ($1,$2)`,
		txID,
		ledger.ReferenceID,
	)
	if err != nil {
		tx.Rollback(ctx)
		fmt.Println("Transaction insert error")
		return err
	}

	// Debit entry
	_, err = tx.Exec(
		ctx,
		`INSERT INTO ledger_entries ( transaction_id, account_id, amount, type)
		 VALUES ($1,$2,$3,$4)`,
		txID,
		ledger.SenderAccount,
		-ledger.Amount,
		"debit",
	)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	// Credit entry
	_, err = tx.Exec(
		ctx,
		`INSERT INTO ledger_entries (transaction_id, account_id, amount, type)
		 VALUES ($1,$2,$3,$4)`,
		txID,
		ledger.ReceiverAccount,
		ledger.Amount,
		"credit",
	)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	// Update balances
	_, err = tx.Exec(
		ctx,
		`UPDATE accounts SET balance = balance - $1 WHERE id = $2`,
		ledger.Amount,
		ledger.SenderAccount,
	)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	_, err = tx.Exec(
		ctx,
		`UPDATE accounts SET balance = balance + $1 WHERE id = $2`,
		ledger.Amount,
		ledger.ReceiverAccount,
	)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
