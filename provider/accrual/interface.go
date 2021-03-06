//go:generate mockgen -source=interface.go -destination=./mock/provider.go -package=providermock
package accrual

import (
	"github.com/vstdy/gophermart/model"
)

type Provider interface {
	// GetOrderAccruals gets order status and accruals.
	GetOrderAccruals(order model.Order) (model.Order, error)
}
