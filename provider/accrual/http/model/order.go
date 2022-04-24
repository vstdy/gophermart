package model

import (
	"github.com/google/uuid"

	"github.com/vstdy/gophermart/model"
)

// Order keeps order data.
type Order struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual"`
}

// ToCanonical converts a accrual model to canonical model.
func (o Order) ToCanonical(userID uuid.UUID) (model.Order, error) {
	obj := model.Order{
		UserID:  userID,
		Number:  o.Order,
		Status:  model.NewOrderStatusFromStr(o.Status),
		Accrual: o.Accrual,
	}

	return obj, nil
}
