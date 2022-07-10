//go:generate mockgen -source=interface.go -destination=./mock/accrual.go -package=accrualmock
package accrual

import (
	"github.com/vstdy/gophermart/model"
)

type Accrual interface {
	// GetOrderAccruals gets order status and accruals.
	GetOrderAccruals(order model.Order) (model.Order, error)
}
