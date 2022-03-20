package schema

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/vstdy0/go-diploma/model"
)

type (
	// Order keeps order data.
	Order struct {
		bun.BaseModel `bun:"orders,alias:o"`
		ID            uuid.UUID `bun:"id,pk,type:uuid"`
		UserID        uuid.UUID `bun:"user_id,type:uuid,notnull"`
		Number        string    `bun:"number,unique,notnull"`
		Status        string    `bun:"status,nullzero,notnull,default:'NEW'"`
		Accrual       int       `bun:"accrual,notnull"`
		UploadedAt    time.Time `bun:"uploaded_at,nullzero,notnull,default:current_timestamp"`
		UpdatedAt     time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
		Updated       bool      `bun:"updated,scanonly"`
	}

	Orders []Order
)

// NewOrderFromCanonical creates a new Order DB object from canonical model.
func NewOrderFromCanonical(obj model.Order) Order {
	return Order{
		UserID:     obj.UserID,
		Number:     obj.Number,
		Status:     obj.Status.String(),
		Accrual:    int(obj.Accrual * 100),
		UploadedAt: obj.UploadedAt,
	}
}

// NewOrdersFromCanonical creates new list of Order DB objects from list of canonical models.
func NewOrdersFromCanonical(objs []model.Order) Orders {
	var orders Orders
	for _, order := range objs {
		orders = append(orders, NewOrderFromCanonical(order))
	}

	return orders
}

// ToCanonical converts a Order DB object to canonical model.
func (o Order) ToCanonical() (model.Order, error) {
	return model.Order{
		UserID:     o.UserID,
		Number:     o.Number,
		Status:     model.NewOrderStatusFromStr(o.Status),
		Accrual:    float32(o.Accrual) / 100,
		UploadedAt: o.UploadedAt,
	}, nil
}

// ToCanonical converts list of Order DB objects to list of canonical models.
func (o Orders) ToCanonical() ([]model.Order, error) {
	objs := make([]model.Order, 0, len(o))
	for _, dbObj := range o {
		obj, err := dbObj.ToCanonical()
		if err != nil {
			return nil, err
		}
		objs = append(objs, obj)
	}

	return objs, nil
}
