package gophermart

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/vstdy/gophermart/model"
	"github.com/vstdy/gophermart/pkg/logging"
	"github.com/vstdy/gophermart/provider/accrual"
	"github.com/vstdy/gophermart/provider/notification"
	"github.com/vstdy/gophermart/service/gophermart"
	"github.com/vstdy/gophermart/storage"
)

const (
	serviceName = "gophermart"
)

var _ gophermart.Service = (*Service)(nil)

type (
	// Service keeps service dependencies.
	Service struct {
		config   Config
		provider provider
		storage  storage.Storage
	}

	// provider keeps provider dependencies.
	provider struct {
		ntf notification.Notification
		ch  chan []model.Transaction

		accrual accrual.Accrual
	}

	// ServiceOption defines functional argument for Service constructor.
	ServiceOption func(*Service) error
)

// WithConfig sets Config.
func WithConfig(config Config) ServiceOption {
	return func(svc *Service) error {
		svc.config = config

		return nil
	}
}

// WithProvider sets provider.
func WithProvider(acc accrual.Accrual, ntf notification.Notification) ServiceOption {
	return func(svc *Service) error {
		svc.provider.accrual = acc
		svc.provider.ntf = ntf
		svc.provider.ch = make(chan []model.Transaction)

		return nil
	}
}

// WithStorage sets Storage.
func WithStorage(st storage.Storage) ServiceOption {
	return func(svc *Service) error {
		svc.storage = st

		return nil
	}
}

// NewService creates a new gophermart service.
func NewService(ctx context.Context, opts ...ServiceOption) (*Service, error) {
	svc := &Service{
		config: NewDefaultConfig(),
	}
	for optIdx, opt := range opts {
		if err := opt(svc); err != nil {
			return nil, fmt.Errorf("applying option [%d]: %w", optIdx, err)
		}
	}

	if svc.storage == nil {
		return nil, fmt.Errorf("storage: nil")
	}

	go svc.orderStatusUpdater(ctx)
	go svc.accrualNotifier(ctx)
	go svc.provider.ntf.ConsumeAccrualNotifications(ctx)

	return svc, nil
}

// Close closes all service dependencies.
func (svc *Service) Close() error {
	if svc.storage == nil {
		return nil
	}
	if err := svc.storage.Close(); err != nil {
		return fmt.Errorf("closing storage: %w", err)
	}

	if err := svc.provider.ntf.Close(); err != nil {
		return fmt.Errorf("closing kafka: %w", err)
	}

	return nil
}

// Logger returns logger with service context.
func (svc *Service) Logger() zerolog.Logger {
	logCtx := log.With().Str(logging.ServiceKey, serviceName)

	return logCtx.Logger()
}
