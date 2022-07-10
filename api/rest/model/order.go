package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/vstdy/gophermart/model"
)

type Order struct {
	UserID     uuid.UUID `json:"-"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float32   `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
}

func NewOrder(userID uuid.UUID, number string) Order {
	return Order{
		UserID: userID,
		Number: number,
	}
}

// NewOrdersFromCanonical creates a new Order object from canonical model.
func NewOrdersFromCanonical(objs []model.Order) []Order {
	var orders []Order
	for _, obj := range objs {
		orders = append(orders, Order{
			UserID:     obj.UserID,
			Number:     obj.Number,
			Status:     obj.Status.String(),
			Accrual:    obj.Accrual,
			UploadedAt: obj.UploadedAt,
		})
	}

	return orders
}

// ToCanonical converts a API model to canonical model.
func (o Order) ToCanonical() model.Order {
	obj := model.Order{
		UserID:     o.UserID,
		Number:     o.Number,
		Status:     model.NewOrderStatusFromStr(o.Status),
		Accrual:    o.Accrual,
		UploadedAt: o.UploadedAt,
	}

	return obj
}

// MarshalJSON implements interface json.Marshaler.
func (o Order) MarshalJSON() ([]byte, error) {
	type OrderAlias Order

	if o.Status != model.OrderStatusProcessed.String() {
		order := struct {
			OrderAlias
			Accrual    int    `json:"accrual,omitempty"`
			UploadedAt string `json:"uploaded_at"`
		}{
			OrderAlias: OrderAlias(o),
			Accrual:    0,
			UploadedAt: o.UploadedAt.Format(time.RFC3339),
		}

		return json.Marshal(order)
	}

	order := struct {
		OrderAlias
		UploadedAt string `json:"uploaded_at"`
	}{
		OrderAlias: OrderAlias(o),
		UploadedAt: o.UploadedAt.Format(time.RFC3339),
	}

	return json.Marshal(order)
}
