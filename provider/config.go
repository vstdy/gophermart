package provider

import (
	"fmt"

	inter "github.com/vstdy/gophermart/provider/accrual"
	"github.com/vstdy/gophermart/provider/accrual/http"
	"github.com/vstdy/gophermart/provider/notification"
	"github.com/vstdy/gophermart/provider/notification/kafka"
)

// Config keeps Provider params.
type Config struct {
	Accrual accrual.Config `mapstructure:"accrual,squash"`
	Kafka   kafka.Config   `mapstructure:"kafka,squash"`
}

// BuildProvider builds Provider dependency.
func (config Config) BuildProvider() (inter.Accrual, notification.Notification, error) {
	acc, err := accrual.NewAccrualProvider(
		config.Accrual.Timeout,
		accrual.WithConfig(config.Accrual),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("building accrual provider: %w", err)
	}

	ntf, err := kafka.NewKafkaProvider(
		kafka.WithConfig(config.Kafka),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("building kafka provider: %w", err)
	}

	return acc, ntf, err
}

// Validate performs a basic validation.
func (config Config) Validate() error {
	if err := config.Accrual.Validate(); err != nil {
		return fmt.Errorf("provider: %w", err)
	}

	if err := config.Kafka.Validate(); err != nil {
		return fmt.Errorf("provider: %w", err)
	}

	return nil
}

// NewDefaultConfig builds a Config with default values.
func NewDefaultConfig() Config {
	return Config{
		Accrual: accrual.NewDefaultConfig(),
		Kafka:   kafka.NewDefaultConfig(),
	}
}
