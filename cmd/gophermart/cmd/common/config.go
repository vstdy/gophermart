package common

import (
	"context"
	"fmt"
	"time"

	"github.com/vstdy0/go-diploma/pkg"
	"github.com/vstdy0/go-diploma/service/gophermart/v1"
	"github.com/vstdy0/go-diploma/storage"
	"github.com/vstdy0/go-diploma/storage/psql"
)

// Config combines sub-configs for all services, storages and providers.
type Config struct {
	Timeout     time.Duration
	RunAddress  string            `mapstructure:"run_address"`
	SecretKey   string            `mapstructure:"secret_key"`
	StorageType string            `mapstructure:"storage_type"`
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

	ctx, ctxCancel := context.WithTimeout(context.Background(), config.Timeout)
	defer ctxCancel()

	if err = st.Migrate(ctx); err != nil {
		return nil, err
	}

	return st, nil
}

// BuildService builds gophermart.Service dependency.
func (config Config) BuildService(storageType string) (*gophermart.Service, error) {
	var st storage.Storage
	var err error

	switch storageType {
	case psqlStorage:
		st, err = config.BuildPsqlStorage()
	default:
		err = pkg.ErrUnsupportedStorageType
	}
	if err != nil {
		return nil, fmt.Errorf("building storage: %w", err)
	}

	svc, err := gophermart.New(
		gophermart.WithConfig(config.Service),
		gophermart.WithStorage(st),
	)
	if err != nil {
		return nil, fmt.Errorf("building service: %w", err)
	}

	return svc, nil
}
