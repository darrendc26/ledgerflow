package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool() *pgxpool.Pool {
	connStr := "postgres://ledgerflow:ledgerflow@postgres:5432/ledgerflow?sslmode=disable"

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	return pool
}
