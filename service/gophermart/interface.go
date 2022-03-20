package gophermart

import (
	"context"
	"io"

	"github.com/google/uuid"

	"github.com/vstdy0/go-diploma/model"
)

type Service interface {
	io.Closer

	// CreateUser creates a new model.User.
	CreateUser(ctx context.Context, obj model.User) (model.User, error)
	// AuthenticateUser verifies the identity of credentials.
	AuthenticateUser(ctx context.Context, obj model.User) (model.User, error)

	// AddOrder adds given order to storage.
	AddOrder(ctx context.Context, obj model.Order) (model.Order, error)
	// GetOrders gets current user orders.
	GetOrders(ctx context.Context, userID uuid.UUID) ([]model.Order, error)

	// GetBalance gets current user balance.
	GetBalance(ctx context.Context, userID uuid.UUID) (float32, float32, error)
	// AddWithdrawal adds withdrawal.
	AddWithdrawal(ctx context.Context, transaction model.Transaction) error
	// GetWithdrawals gets current user withdrawals.
	GetWithdrawals(ctx context.Context, userID uuid.UUID) ([]model.Transaction, error)
}
