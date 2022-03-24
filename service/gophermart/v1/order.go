package gophermart

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/vstdy0/go-diploma/model"
	"github.com/vstdy0/go-diploma/service/gophermart/v1/validator"
)

// AddOrder adds given order object to storage.
func (svc *Service) AddOrder(ctx context.Context, obj model.Order) (model.Order, error) {
	if err := validator.ValidateOrderNumber(obj.Number); err != nil {
		return model.Order{}, err
	}

	addedObj, err := svc.storage.AddOrder(ctx, obj)
	if err != nil {
		return addedObj, err
	}

	return addedObj, nil
}

// GetOrders gets current user orders.
func (svc *Service) GetOrders(ctx context.Context, userID uuid.UUID) ([]model.Order, error) {
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

		var orders []model.Order
		var transactions []model.Transaction
		for _, obj := range objs {
			order, err := svc.provider.GetOrderAccruals(obj)
			if err != nil {
				return fmt.Errorf("accrual provider: %w", err)
			}

			if order.Status.Validate() != nil {
				continue
			}

			if order.Status == model.OrderStatusProcessed && order.Accrual > 0 {
				transactions = append(transactions, model.NewTransaction(order))
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
