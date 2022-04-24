package gophermart

import (
	"context"

	"github.com/google/uuid"

	"github.com/vstdy/gophermart/model"
	"github.com/vstdy/gophermart/service/gophermart/v1/validator"
)

// GetBalance gets current user balance.
func (svc *Service) GetBalance(ctx context.Context, userID uuid.UUID) (float32, float32, error) {
	current, used, err := svc.storage.GetBalance(ctx, userID)
	if err != nil {
		return 0, 0, err
	}

	return current, used, nil
}

// AddWithdrawal adds withdrawal.
func (svc *Service) AddWithdrawal(ctx context.Context, transaction model.Transaction) error {
	if err := validator.ValidateOrderNumber(transaction.Order); err != nil {
		return err
	}

	err := svc.storage.AddWithdrawal(ctx, transaction)
	if err != nil {
		return err
	}

	return nil
}

// GetWithdrawals gets current user withdrawals.
func (svc *Service) GetWithdrawals(ctx context.Context, userID uuid.UUID) ([]model.Transaction, error) {
	objs, err := svc.storage.GetWithdrawals(ctx, userID)
	if err != nil {
		return nil, err
	}

	return objs, nil
}
