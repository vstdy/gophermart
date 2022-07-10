package gophermart

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	canonical "github.com/vstdy/gophermart/model"
	"github.com/vstdy/gophermart/service/gophermart/v1/validator"
)

// AddOrder adds given order object to storage.
func (svc *Service) AddOrder(ctx context.Context, obj canonical.Order) (canonical.Order, error) {
	if err := validator.ValidateOrderNumber(obj.Number); err != nil {
		return canonical.Order{}, err
	}

	addedObj, err := svc.storage.AddOrder(ctx, obj)
	if err != nil {
		return addedObj, err
	}

	return addedObj, nil
}

// GetOrders gets current user orders.
func (svc *Service) GetOrders(ctx context.Context, userID uuid.UUID) ([]canonical.Order, error) {
	objs, err := svc.storage.GetOrders(ctx, userID)
	if err != nil {
		return objs, err
	}

	return objs, nil
}

// orderStatusUpdater updates orders objects status.
func (svc *Service) orderStatusUpdater(ctx context.Context) {
	update := func() error {
		updCtx, cancel := context.WithTimeout(context.Background(), svc.config.UpdaterTimeout)
		defer cancel()

		objs, err := svc.storage.GetStatusNewOrders(updCtx)
		if err != nil {
			return fmt.Errorf("get orders objects: %w", err)
		}

		var orders []canonical.Order
		var transactions []canonical.Transaction
		for _, obj := range objs {
			order, err := svc.provider.accrual.GetOrderAccruals(obj)
			if err != nil {
				return fmt.Errorf("accrual provider: %w", err)
			}

			if order.Status.Validate() != nil {
				continue
			}

			if order.Status == canonical.OrderStatusProcessed && order.Accrual > 0 {
				transactions = append(transactions, canonical.NewTransaction(order))
			}

			orders = append(orders, order)
		}

		if len(orders) > 0 {
			updCtx, cancel = context.WithTimeout(context.Background(), svc.config.UpdaterTimeout)
			defer cancel()

			if err = svc.storage.UpdateOrders(updCtx, orders); err != nil {
				return fmt.Errorf("update orders objects: %w", err)
			}
		}

		if len(transactions) > 0 {
			updCtx, cancel = context.WithTimeout(context.Background(), svc.config.UpdaterTimeout)
			defer cancel()

			if err = svc.storage.AddAccruals(updCtx, transactions); err != nil {
				return fmt.Errorf("add accruals: %w", err)
			}

			go func(tr []canonical.Transaction) {
				svc.provider.ch <- tr
			}(transactions)
		}

		return nil
	}

	ticker := time.NewTicker(svc.config.StatusCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("orderStatusUpdater closed")
			return
		case <-ticker.C:
			if err := update(); err != nil {
				log.Warn().Err(err).Msg("orderStatusUpdater:")
			}
		}
	}
}
