package model

type Payment struct {
	ID       string
	UserID   string
	Amount   int64
	Currency string
	Status   string
}
