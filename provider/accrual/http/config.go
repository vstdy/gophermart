package accrual

import (
	"fmt"
	"time"
)

// Config keeps Accrual params.
type Config struct {
	Timeout           time.Duration `mapstructure:"timeout"`
	AccrualSysAddress string        `mapstructure:"accrual_system_address"`
}

// Validate performs a basic validation.
func (config Config) Validate() error {
	if config.Timeout < time.Second {
		return fmt.Errorf("timeout field: must be GTE 1")
	}

	if config.AccrualSysAddress == "" {
		return fmt.Errorf("accrual_system_address field: empty")
	}

	return nil
}

// NewDefaultConfig builds a Config with default values.
func NewDefaultConfig() Config {
	return Config{
		Timeout:           5 * time.Second,
		AccrualSysAddress: "http://127.0.0.1:8081",
	}
}
