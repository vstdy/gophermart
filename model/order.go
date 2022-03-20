package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Order keeps order data.
type Order struct {
	UserID     uuid.UUID   `json:"user_id"`
	Number     string      `json:"order"`
	Status     OrderStatus `json:"status"`
	Accrual    float32     `json:"accrual"`
	UploadedAt time.Time   `json:"uploaded_at"`
}

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

// NewOrderStatusFromStr returns OrderStatus by its str representation (might be invalid).
func NewOrderStatusFromStr(o string) OrderStatus {
	return OrderStatus(o)
}

// String implements fmt.Stringer interface.
func (o OrderStatus) String() string {
	return string(o)
}

// Validate performs enum validation.
func (o OrderStatus) Validate() error {
	switch o {
	case OrderStatusNew, OrderStatusProcessing, OrderStatusInvalid, OrderStatusProcessed:
		return nil
	default:
		return fmt.Errorf("unknown OrderStatus: %s", o)
	}
}
