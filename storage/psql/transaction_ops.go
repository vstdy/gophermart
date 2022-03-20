package psql

import (
	"context"

	"github.com/google/uuid"

	"github.com/vstdy0/go-diploma/model"
	"github.com/vstdy0/go-diploma/pkg"
	"github.com/vstdy0/go-diploma/storage/psql/schema"
)

const transactionTableName = "transaction"

// GetBalance gets current user balance.
func (st *Storage) GetBalance(ctx context.Context, userID uuid.UUID) (float32, float32, error) {
	var dbObj schema.Transaction
	var dbCurrent int
	var dbUsed int

	err := st.db.NewSelect().
		Model(&dbObj).
		ColumnExpr("sum(accrual) - sum(withdrawal) AS current, sum(withdrawal) AS used").
		Where("user_id = ?", userID).
		Scan(ctx, &dbCurrent, &dbUsed)
	if err != nil {
		return 0, 0, err
	}

	current := float32(dbCurrent) / 100
	used := float32(dbUsed) / 100

	return current, used, nil
}

// AddAccruals adds accruals.
func (st *Storage) AddAccruals(ctx context.Context, objs []model.Transaction) error {
	dbObjs := schema.NewTransactionsFromCanonical(objs)

	_, err := st.db.NewInsert().
		Model(&dbObjs).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// AddWithdrawal adds withdrawal.
func (st *Storage) AddWithdrawal(ctx context.Context, obj model.Transaction) error {
	dbObj := schema.NewTransactionFromCanonical(obj)
	var enough bool

	err := st.db.NewSelect().
		Model(&dbObj).
		ColumnExpr("sum(accrual) - sum(withdrawal) > ? AS enough", dbObj.Withdrawal).
		Where("user_id = ?", dbObj.UserID).
		Scan(ctx, &enough)
	if err != nil {
		return err
	}

	if !enough {
		return pkg.ErrNonSufficientFunds
	}

	_, err = st.db.NewInsert().
		Model(&dbObj).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// GetWithdrawals gets current user withdrawals.
func (st *Storage) GetWithdrawals(ctx context.Context, userID uuid.UUID) ([]model.Transaction, error) {
	var dbObjs schema.Transactions

	err := st.db.NewSelect().
		Model(&dbObjs).
		Where("user_id = ?", userID).
		Where("withdrawal > 0").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	if dbObjs == nil {
		return nil, nil
	}

	objs, err := dbObjs.ToCanonical()
	if err != nil {
		return nil, err
	}

	return objs, nil
}
