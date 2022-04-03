package schema

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/vstdy0/go-diploma/model"
)

// Transaction keeps order data.
type (
	Transaction struct {
		bun.BaseModel `bun:"transactions,alias:t"`
		ID            uuid.UUID `bun:"id,pk,type:uuid"`
		UserID        uuid.UUID `bun:"user_id,type:uuid,notnull"`
		Order         string    `bun:"order,unique,notnull"`
		Accrual       int       `bun:"accrual,notnull"`
		Withdrawal    int       `bun:"withdrawal,notnull"`
		ProcessedAt   time.Time `bun:"processed_at,nullzero,notnull,default:current_timestamp"`
	}

	Transactions []Transaction
)

// NewTransactionFromCanonical creates a new Transaction DB object from canonical model.
func NewTransactionFromCanonical(obj model.Transaction) Transaction {
	return Transaction{
		UserID:     obj.UserID,
		Order:      obj.Order,
		Accrual:    int(obj.Accrual * 100),
		Withdrawal: int(obj.Withdrawal * 100),
	}
}

// NewTransactionsFromCanonical creates new list of Transaction DB objects from list of canonical models.
func NewTransactionsFromCanonical(objs []model.Transaction) []Transaction {
	var transactions []Transaction
	for _, transaction := range objs {
		transactions = append(transactions, NewTransactionFromCanonical(transaction))
	}

	return transactions
}

// ToCanonical converts a Order DB object to canonical model.
func (o Transaction) ToCanonical() (model.Transaction, error) {
	return model.Transaction{
		UserID:      o.UserID,
		Order:       o.Order,
		Accrual:     float32(o.Accrual) / 100,
		Withdrawal:  float32(o.Withdrawal) / 100,
		ProcessedAt: o.ProcessedAt,
	}, nil
}

// ToCanonical converts list of Order DB objects to list of canonical models.
func (o Transactions) ToCanonical() ([]model.Transaction, error) {
	objs := make([]model.Transaction, 0, len(o))
	for _, dbObj := range o {
		obj, err := dbObj.ToCanonical()
		if err != nil {
			return nil, err
		}
		objs = append(objs, obj)
	}

	return objs, nil
}
