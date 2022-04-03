package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/vstdy0/go-diploma/model"
)

type BalanceResponse struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

type AddWithdrawalBody struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}

// ToCanonical converts a API model to canonical model.
func (w AddWithdrawalBody) ToCanonical(userID uuid.UUID) model.Transaction {
	obj := model.Transaction{
		UserID:     userID,
		Order:      w.Order,
		Withdrawal: w.Sum,
	}

	return obj
}

type (
	GetWithdrawal struct {
		Order       string    `json:"order"`
		Sum         float32   `json:"sum"`
		ProcessedAt time.Time `json:"processed_at"`
	}

	GetWithdrawals []GetWithdrawal
)

// NewGetWithdrawalFromCanonical creates a new Transaction DB object from canonical model.
func NewGetWithdrawalFromCanonical(obj model.Transaction) GetWithdrawal {
	return GetWithdrawal{
		Order:       obj.Order,
		Sum:         obj.Withdrawal,
		ProcessedAt: obj.ProcessedAt,
	}
}

// NewGetWithdrawalsFromCanonical creates new list of Transaction DB objects from list of canonical models.
func NewGetWithdrawalsFromCanonical(objs []model.Transaction) GetWithdrawals {
	var getWithdrawals GetWithdrawals
	for _, transaction := range objs {
		getWithdrawals = append(getWithdrawals, NewGetWithdrawalFromCanonical(transaction))
	}

	return getWithdrawals
}

// MarshalJSON implements interface json.Marshaler.
func (w GetWithdrawal) MarshalJSON() ([]byte, error) {
	type GetWithdrawalAlias GetWithdrawal

	getWithdrawal := struct {
		GetWithdrawalAlias
		ProcessedAt string `json:"processed_at"`
	}{
		GetWithdrawalAlias: GetWithdrawalAlias(w),
		ProcessedAt:        w.ProcessedAt.Format(time.RFC3339),
	}

	return json.Marshal(getWithdrawal)
}
