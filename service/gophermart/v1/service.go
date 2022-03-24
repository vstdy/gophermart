package gophermart

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/vstdy0/go-diploma/pkg/logging"
	"github.com/vstdy0/go-diploma/provider/accrual"
	"github.com/vstdy0/go-diploma/service/gophermart"
	"github.com/vstdy0/go-diploma/storage"
)

const (
	serviceName = "gophermart"
)

var _ gophermart.Service = (*Service)(nil)

type (
	// Service keeps service dependencies.
	Service struct {
		config   Config
		provider accrual.Provider
		storage  storage.Storage
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

// WithProvider sets Provider.
func WithProvider(p accrual.Provider) ServiceOption {
	return func(svc *Service) error {
		svc.provider = p

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

// New creates a new gophermart service.
func New(ctx context.Context, opts ...ServiceOption) (*Service, error) {
	svc := &Service{
		config: NewDefaultConfig(),
	}
	for optIdx, opt := range opts {
		if err := opt(svc); err != nil {
			return nil, fmt.Errorf("applying option [%d]: %w", optIdx, err)
		}
	}

	if err := svc.config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation: %w", err)
	}

	if svc.storage == nil {
		return nil, fmt.Errorf("storage: nil")
	}

	go svc.orderStatusUpdater(ctx)

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

	return nil
}

// Logger returns logger with service context.
func (svc *Service) Logger() zerolog.Logger {
	logCtx := log.With().Str(logging.ServiceKey, serviceName)

	return logCtx.Logger()
}
