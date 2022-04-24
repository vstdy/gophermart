//go:generate mockgen -source=interface.go -destination=./mock/storage.go -package=storagemock
package storage

import (
	"context"
	"io"

	"github.com/google/uuid"

	"github.com/vstdy/gophermart/model"
)

type Storage interface {
	io.Closer

	// CreateUser adds given objects to storage.
	CreateUser(ctx context.Context, obj model.User) (model.User, error)
	// AuthenticateUser verifies the identity of credentials.
	AuthenticateUser(ctx context.Context, obj model.User) (model.User, error)

	// AddOrder adds given order to storage.
	AddOrder(ctx context.Context, obj model.Order) (model.Order, error)
	// GetStatusNewOrders gets orders with status new.
	GetStatusNewOrders(ctx context.Context) ([]model.Order, error)
	// UpdateOrders updates given orders.
	UpdateOrders(ctx context.Context, objs []model.Order) error
	// GetOrders gets current user orders.
	GetOrders(ctx context.Context, userID uuid.UUID) ([]model.Order, error)

	// GetBalance gets current user balance.
	GetBalance(ctx context.Context, userID uuid.UUID) (float32, float32, error)
	// AddAccruals updates given orders.
	AddAccruals(ctx context.Context, objs []model.Transaction) error
	// AddWithdrawal adds withdrawal.
	AddWithdrawal(ctx context.Context, transaction model.Transaction) error
	// GetWithdrawals gets current user withdrawals.
	GetWithdrawals(ctx context.Context, userID uuid.UUID) ([]model.Transaction, error)
}
