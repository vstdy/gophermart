package common

import (
	"context"
	"fmt"

	"github.com/vstdy/gophermart/api"
	"github.com/vstdy/gophermart/provider"
	inter "github.com/vstdy/gophermart/service/gophermart"
	"github.com/vstdy/gophermart/service/gophermart/v1"
	storage "github.com/vstdy/gophermart/storage/common"
)

// Config combines sub-configs for all services, storages and providers.
type Config struct {
	API      api.Config        `mapstructure:"server,squash"`
	Provider provider.Config   `mapstructure:"provider,squash"`
	Service  gophermart.Config `mapstructure:"service,squash"`
	Storage  storage.Config    `mapstructure:"psql_storage,squash"`
}

// Validate performs a basic validation.
func (config Config) Validate() error {
	if err := config.API.Validate(); err != nil {
		return fmt.Errorf("API: %w", err)
	}

	if err := config.Provider.Validate(); err != nil {
		return fmt.Errorf("provider: %w", err)
	}

	if err := config.Service.Validate(); err != nil {
		return fmt.Errorf("service: %w", err)
	}

	if err := config.Storage.Validate(); err != nil {
		return fmt.Errorf("storage: %w", err)
	}

	return nil
}

// BuildDefaultConfig builds a Config with default values.
func BuildDefaultConfig() Config {
	return Config{
		API:      api.NewDefaultConfig(),
		Provider: provider.NewDefaultConfig(),
		Service:  gophermart.NewDefaultConfig(),
		Storage:  storage.NewDefaultConfig(),
	}
}

// BuildService builds gophermart.Service dependency.
func (config Config) BuildService(ctx context.Context) (inter.Service, error) {
	acc, ntf, err := config.Provider.BuildProvider()
	if err != nil {
		return nil, fmt.Errorf("building provider: %w", err)
	}

	st, err := config.Storage.BuildStorage()
	if err != nil {
		return nil, fmt.Errorf("building storage: %w", err)
	}

	svc, err := gophermart.NewService(
		ctx,
		gophermart.WithConfig(config.Service),
		gophermart.WithProvider(acc, ntf),
		gophermart.WithStorage(st),
	)
	if err != nil {
		return nil, fmt.Errorf("building service: %w", err)
	}

	return svc, nil
}
