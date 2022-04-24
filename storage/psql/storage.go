package psql

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"
	"sync"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"

	inter "github.com/vstdy/gophermart/storage"
	"github.com/vstdy/gophermart/storage/psql/migrations"
	"github.com/vstdy/gophermart/storage/psql/schema"
)

const (
	serviceName = "psql"

	dbTableLoggingKey     = "db-table"
	dbOperationLoggingKey = "db-operation"
)

var _ inter.Storage = (*Storage)(nil)

type (
	// Storage keeps psql storage dependencies.
	Storage struct {
		sync.RWMutex

		config Config
		db     *bun.DB
	}

	// StorageOption defines functional argument for Storage constructor.
	StorageOption func(st *Storage) error
)

// WithConfig overrides default Storage config.
func WithConfig(config Config) StorageOption {
	return func(st *Storage) error {
		st.config = config

		return nil
	}
}

// New creates a new psql Storage with custom options.
func New(opts ...StorageOption) (*Storage, error) {
	st := &Storage{
		config: NewDefaultConfig(),
	}
	for optIdx, opt := range opts {
		if err := opt(st); err != nil {
			return nil, fmt.Errorf("applying option [%d]: %w", optIdx, err)
		}
	}

	if err := st.config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation: %w", err)
	}

	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(st.config.URI)))

	maxOpenConnections := 4 * runtime.GOMAXPROCS(0)

	st.db = bun.NewDB(sqlDB, pgdialect.New())
	st.db.AddQueryHook(newQueryHook(st))
	st.db.SetMaxOpenConns(maxOpenConnections)
	st.db.SetMaxIdleConns(maxOpenConnections)
	st.db.RegisterModel(
		(*schema.User)(nil),
	)

	if err := st.db.Ping(); err != nil {
		return nil, fmt.Errorf("ping for URI (%s) failed: %w", st.config.URI, err)
	}

	return st, nil
}

// Close closes DB connection.
func (st *Storage) Close() error {
	if st.db == nil {
		return nil
	}

	return st.db.Close()
}

// Migrate performs DB migrations.
func (st *Storage) Migrate(ctx context.Context) error {
	logger := st.Logger(withOperation("migration"))

	ms, err := migrations.GetMigrations()
	if err != nil {
		return err
	}

	migration := migrate.NewMigrator(st.db, ms)

	if err = migration.Init(ctx); err != nil {
		return fmt.Errorf("initialising migration: %w", err)
	}

	res, err := migration.Migrate(ctx)
	if err != nil {
		return fmt.Errorf("performing migration: %w", err)
	}

	logger.Info().Msgf("Migration applied: %s", res.Migrations.LastGroup().String())

	return nil
}
