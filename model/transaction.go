package model

import (
	"time"

	"github.com/google/uuid"
)

// Transaction keeps order data.
type Transaction struct {
	UserID      uuid.UUID
	Order       string
	Accrual     float32
	Withdrawal  float32
	ProcessedAt time.Time
}

// NewTransaction creates a new Transaction model from Order model.
func NewTransaction(obj Order) Transaction {
	return Transaction{
		UserID:  obj.UserID,
		Order:   obj.Number,
		Accrual: obj.Accrual,
	}
}
