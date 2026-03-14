package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) CreateAccount(userID string, currency string) (string, error) {

	var accountID string

	query := `
	INSERT INTO accounts (user_id, balance, currency)
	VALUES ($1, 0, $2)
	RETURNING id
	`

	err := r.db.QueryRow(context.Background(), query, userID, currency).Scan(&accountID)

	return accountID, err
}
