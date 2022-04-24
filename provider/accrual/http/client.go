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

var _ accrual.Provider = (*Provider)(nil)

// WithConfig sets Config.
func WithConfig(config Config) ProviderOption {
	return func(svc *Provider) error {
		svc.config = config

		return nil
	}
}

// Provider keeps accrual service configuration.
type (
	Provider struct {
		config Config
		client http.Client
	}

	// ProviderOption defines functional argument for Provider constructor.
	ProviderOption func(*Provider) error
)

// NewProvider returns a new Provider instance.
func NewProvider(timeout time.Duration, opts ...ProviderOption) (*Provider, error) {
	prv := &Provider{
		config: NewDefaultConfig(),
	}
	for optIdx, opt := range opts {
		if err := opt(prv); err != nil {
			return nil, fmt.Errorf("applying option [%d]: %w", optIdx, err)
		}
	}

	if err := prv.config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation: %w", err)
	}

	prv.client = http.Client{Timeout: timeout}
	transport := &http.Transport{}
	transport.MaxIdleConns = 1
	prv.client.Transport = transport

	return prv, nil
}

// GetOrderAccruals implements the accrual.Provider interface.
func (p Provider) GetOrderAccruals(obj canonical.Order) (canonical.Order, error) {
	url := fmt.Sprintf("%s/api/orders/%s", p.config.AccrualSysAddress, obj.Number)
	r, err := p.client.Get(url)
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
