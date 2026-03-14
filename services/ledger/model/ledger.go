package model

type Ledger struct {
	SenderAccount   string
	ReceiverAccount string
	Amount          int64
	ReferenceID     string
	TransactionID   string
	Status          string
}
