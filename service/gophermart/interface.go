package gophermart

import (
	"context"
	"io"

	"github.com/google/uuid"

	canonical "github.com/vstdy/gophermart/model"
)

type Service interface {
	io.Closer

	// CreateUser creates a new model.User.
	CreateUser(ctx context.Context, obj canonical.User) (canonical.User, error)
	// AuthenticateUser verifies the identity of credentials.
	AuthenticateUser(ctx context.Context, obj canonical.User) (canonical.User, error)

	// AddOrder adds given order to storage.
	AddOrder(ctx context.Context, obj canonical.Order) (canonical.Order, error)
	// GetOrders gets current user orders.
	GetOrders(ctx context.Context, userID uuid.UUID) ([]canonical.Order, error)

	// GetBalance gets current user balance.
	GetBalance(ctx context.Context, userID uuid.UUID) (float32, float32, error)
	// AddWithdrawal adds withdrawal.
	AddWithdrawal(ctx context.Context, transaction canonical.Transaction) error
	// GetWithdrawals gets current user withdrawals.
	GetWithdrawals(ctx context.Context, userID uuid.UUID) ([]canonical.Transaction, error)

	// GetAccrualNotificationsChan returns accrual notifications channel.
	GetAccrualNotificationsChan() chan canonical.Transaction
}
