package accrual

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	canonical "github.com/vstdy/gophermart/model"
	"github.com/vstdy/gophermart/provider/accrual"
	"github.com/vstdy/gophermart/provider/accrual/http/model"
)

var _ accrual.Accrual = (*Accrual)(nil)

// Accrual keeps accrual service configuration.
type (
	Accrual struct {
		config Config
		client http.Client
	}

	// AccrualOption defines functional argument for Accrual constructor.
	AccrualOption func(*Accrual) error
)

// WithConfig sets Config.
func WithConfig(config Config) AccrualOption {
	return func(svc *Accrual) error {
		svc.config = config

		return nil
	}
}

// NewAccrualProvider returns a new Accrual instance.
func NewAccrualProvider(timeout time.Duration, opts ...AccrualOption) (*Accrual, error) {
	acr := &Accrual{
		config: NewDefaultConfig(),
	}
	for optIdx, opt := range opts {
		if err := opt(acr); err != nil {
			return nil, fmt.Errorf("applying option [%d]: %w", optIdx, err)
		}
	}

	acr.client = http.Client{Timeout: timeout}
	transport := &http.Transport{}
	transport.MaxIdleConns = 1
	acr.client.Transport = transport

	return acr, nil
}

// GetOrderAccruals implements the Accrual interface.
func (a *Accrual) GetOrderAccruals(obj canonical.Order) (canonical.Order, error) {
	url := fmt.Sprintf("%s/api/orders/%s", a.config.AccrualSysAddress, obj.Number)
	r, err := a.client.Get(url)
	if err != nil {
		return canonical.Order{}, fmt.Errorf("retrieving order object: %w", err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return canonical.Order{}, nil
	}

	var order model.Order

	if err = json.NewDecoder(r.Body).Decode(&order); err != nil {
		return canonical.Order{}, fmt.Errorf("decoding order object: %w", err)
	}

	can, err := order.ToCanonical(obj.UserID)
	if err != nil {
		return canonical.Order{}, fmt.Errorf("converting to canonical: %w", err)
	}

	return can, nil
}
