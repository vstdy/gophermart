package model

import (
	"github.com/google/uuid"

	canonical "github.com/vstdy/gophermart/model"
)

// AccrualNotification keeps order data.
type AccrualNotification struct {
	UserID  uuid.UUID `json:"user_id"`
	Order   string    `json:"order"`
	Accrual float32   `json:"accrual"`
}

// NewAccrualNotificationFromCanonical converts canonical model to notification model.
func NewAccrualNotificationFromCanonical(obj canonical.Transaction) AccrualNotification {
	return AccrualNotification{
		UserID:  obj.UserID,
		Order:   obj.Order,
		Accrual: obj.Accrual,
	}
}

// NewAccrualNotificationsFromCanonical converts canonical model to notification model.
func NewAccrualNotificationsFromCanonical(objs []canonical.Transaction) []AccrualNotification {
	var notifications []AccrualNotification
	for _, obj := range objs {
		notifications = append(notifications, NewAccrualNotificationFromCanonical(obj))
	}

	return notifications
}

// ToCanonical converts a AccrualNotification object to canonical model.
func (an AccrualNotification) ToCanonical() canonical.Transaction {
	return canonical.Transaction{
		UserID:  an.UserID,
		Order:   an.Order,
		Accrual: an.Accrual,
	}
}
