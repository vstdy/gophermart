package common

import (
	"fmt"
	"time"

	"github.com/vstdy/gophermart/pkg"
	inter "github.com/vstdy/gophermart/storage"
	"github.com/vstdy/gophermart/storage/pgbun"
)

const (
	pgBun = "pgbun"
)

// Config keeps api params.
type Config struct {
	Timeout     time.Duration `mapstructure:"timeout"`
	StorageType string        `mapstructure:"storage_type"`
	PgBun       pgbun.Config  `mapstructure:"pgbun,squash"`
}

// BuildStorage builds Storage dependency.
func (config Config) BuildStorage() (inter.Storage, error) {
	var st inter.Storage
	var err error

	switch config.StorageType {
	case pgBun:
		st, err = pgbun.New(pgbun.WithConfig(config.PgBun))
	default:
		err = pkg.ErrUnsupportedStorageType
	}

	return st, err
}

// Validate performs a basic validation.
func (config Config) Validate() error {
	if config.Timeout < time.Second {
		return fmt.Errorf("timeout field: must be GTE 1")
	}

	if config.StorageType == "" {
		return fmt.Errorf("storage_type field: empty")
	}

	if err := config.PgBun.Validate(); err != nil {
		return fmt.Errorf("pgbun: %w", err)
	}

	return nil
}

// NewDefaultConfig builds a Config with default values.
func NewDefaultConfig() Config {
	return Config{
		Timeout:     5 * time.Second,
		StorageType: "pgbun",
		PgBun:       pgbun.NewDefaultConfig(),
	}
}
