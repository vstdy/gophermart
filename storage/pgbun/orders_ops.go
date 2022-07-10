package pgbun

import (
	"context"

	"github.com/google/uuid"

	"github.com/vstdy/gophermart/model"
	"github.com/vstdy/gophermart/pkg"
	"github.com/vstdy/gophermart/storage/pgbun/schema"
)

const orderTableName = "order"

// AddOrder adds given order to storage.
func (st *Storage) AddOrder(ctx context.Context, obj model.Order) (model.Order, error) {
	dbObj := schema.NewOrderFromCanonical(obj)

	_, err := st.db.NewInsert().
		Model(&dbObj).
		On("CONFLICT (\"number\") DO UPDATE").
		Set("updated_at=NOW()").
		Returning("*, uploaded_at <> updated_at AS updated").
		Exec(ctx)
	if err != nil {
		return model.Order{}, err
	}

	addedObj, err := dbObj.ToCanonical()
	if err != nil {
		return model.Order{}, err
	}

	if dbObj.Updated {
		return addedObj, pkg.ErrAlreadyExists
	}

	return addedObj, nil
}

// UpdateOrders updates given orders.
func (st *Storage) UpdateOrders(ctx context.Context, objs []model.Order) error {
	dbObjs := schema.NewOrdersFromCanonical(objs)
	values := st.db.NewValues(&dbObjs)

	_, err := st.db.NewUpdate().
		With("_data", values).
		Model(&dbObjs).
		TableExpr("_data").
		Set("status = _data.status").
		Set("accrual = _data.accrual").
		Where("o.number = _data.number").
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// GetStatusNewOrders gets orders with status NEW.
func (st *Storage) GetStatusNewOrders(ctx context.Context) ([]model.Order, error) {
	var dbObjs schema.Orders

	err := st.db.NewSelect().
		Model(&dbObjs).
		Where("status = 'NEW'").
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

// GetOrders gets current user orders.
func (st *Storage) GetOrders(ctx context.Context, userID uuid.UUID) ([]model.Order, error) {
	var dbObjs schema.Orders

	err := st.db.NewSelect().
		Model(&dbObjs).
		Where("user_id = ?", userID).
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
