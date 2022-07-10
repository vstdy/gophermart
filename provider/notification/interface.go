package notification

import (
	"context"
	"io"

	canonical "github.com/vstdy/gophermart/model"
)

type Notification interface {
	io.Closer

	// ProduceAccrualNotifications sends notifications about accruals.
	ProduceAccrualNotifications(transactions []canonical.Transaction) error
	// ConsumeAccrualNotifications receives notifications about accruals.
	ConsumeAccrualNotifications(ctx context.Context)
	// GetAccrualNotificationsChan returns accrual notifications channel.
	GetAccrualNotificationsChan() chan canonical.Transaction
}
