package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Order keeps order data.
type Order struct {
	UserID     uuid.UUID
	Number     string
	Status     OrderStatus
	Accrual    float32
	UploadedAt time.Time
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
