package gophermart

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
func (svc *Service) orderStatusUpdater() {
	client := http.Client{Timeout: svc.config.UpdaterTimeout}
	transport := &http.Transport{}
	transport.MaxIdleConns = 1
	client.Transport = transport

	update := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), svc.config.UpdaterTimeout)
		defer cancel()

		objs, err := svc.storage.GetStatusNewOrders(ctx)
		if err != nil {
			return fmt.Errorf("get orders objects: %w", err)
		}

		var (
			order        model.Order
			orders       []model.Order
			transactions []model.Transaction
		)
		for _, obj := range objs {
			url := fmt.Sprintf("%s/api/orders/%s", svc.config.AccrualSysAddress, obj.Number)
			r, err := client.Get(url)
			if err != nil {
				return fmt.Errorf("retrieve order object: %w", err)
			}
			defer r.Body.Close()

			if r.StatusCode == http.StatusOK {
				if err = json.NewDecoder(r.Body).Decode(&order); err != nil {
					return fmt.Errorf("decode order object: %w", err)
				}
				order.UserID = obj.UserID

				if order.Status.Validate() != nil {
					continue
				}

				if order.Status == model.OrderStatusProcessed && order.Accrual > 0 {
					transactions = append(transactions, model.NewTransaction(order))
				}

				orders = append(orders, order)
			}
		}

		if orders != nil {
			ctx, cancel = context.WithTimeout(context.Background(), svc.config.UpdaterTimeout)
			defer cancel()

			if err = svc.storage.UpdateOrders(ctx, orders); err != nil {
				return fmt.Errorf("update orders objects: %w", err)
			}
		}

		if transactions != nil {
			ctx, cancel = context.WithTimeout(context.Background(), svc.config.UpdaterTimeout)
			defer cancel()

			if err = svc.storage.AddAccruals(ctx, transactions); err != nil {
				return fmt.Errorf("add accruals: %w", err)
			}
		}

		return nil
	}

	t := time.NewTicker(svc.config.StatusCheckInterval)
	for range t.C {
		if err := update(); err != nil {
			log.Warn().Err(err).Msg("orderStatusUpdater:")
		}
	}
}
