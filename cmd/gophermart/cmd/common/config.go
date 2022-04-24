package common

import (
	"context"
	"fmt"
	"time"

	"github.com/vstdy/gophermart/pkg"
	"github.com/vstdy/gophermart/provider/accrual/http"
	"github.com/vstdy/gophermart/service/gophermart/v1"
	"github.com/vstdy/gophermart/storage"
	"github.com/vstdy/gophermart/storage/psql"
)

// Config combines sub-configs for all services, storages and providers.
type Config struct {
	Timeout     time.Duration     `mapstructure:"timeout"`
	RunAddress  string            `mapstructure:"run_address"`
	SecretKey   string            `mapstructure:"secret_key"`
	StorageType string            `mapstructure:"storage_type"`
	Provider    accrual.Config    `mapstructure:"provider,squash"`
	Service     gophermart.Config `mapstructure:"service,squash"`
	PSQLStorage psql.Config       `mapstructure:"psql_storage,squash"`
}

const (
	psqlStorage = "psql"
)

// BuildDefaultConfig builds a Config with default values.
func BuildDefaultConfig() Config {
	return Config{
		Timeout:     5 * time.Second,
		RunAddress:  "0.0.0.0:8080",
		SecretKey:   "secret_key",
		StorageType: psqlStorage,
		Provider:    accrual.NewDefaultConfig(),
		Service:     gophermart.NewDefaultConfig(),
		PSQLStorage: psql.NewDefaultConfig(),
	}
}

// BuildPsqlStorage builds psql.Storage dependency.
func (config Config) BuildPsqlStorage() (*psql.Storage, error) {
	st, err := psql.New(
		psql.WithConfig(config.PSQLStorage),
	)
	if err != nil {
		return nil, fmt.Errorf("building psql storage: %w", err)
	}

	return st, nil
}

// BuildService builds gophermart.Service dependency.
func (config Config) BuildService(ctx context.Context) (*gophermart.Service, error) {
	var st storage.Storage
	var err error

	prv, err := accrual.NewProvider(
		config.Timeout,
		accrual.WithConfig(config.Provider),
	)
	if err != nil {
		return nil, fmt.Errorf("building provider: %w", err)
	}

	switch config.StorageType {
	case psqlStorage:
		st, err = config.BuildPsqlStorage()
	default:
		err = pkg.ErrUnsupportedStorageType
	}
	if err != nil {
		return nil, fmt.Errorf("building storage: %w", err)
	}

	svc, err := gophermart.New(
		ctx,
		gophermart.WithConfig(config.Service),
		gophermart.WithProvider(prv),
		gophermart.WithStorage(st),
	)
	if err != nil {
		return nil, fmt.Errorf("building service: %w", err)
	}

	return svc, nil
}
