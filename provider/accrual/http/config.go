package accrual

import (
	"fmt"
)

// Config keeps Provider params.
type Config struct {
	AccrualSysAddress string `mapstructure:"accrual_system_address"`
}

// Validate performs a basic validation.
func (config Config) Validate() error {
	if config.AccrualSysAddress == "" {
		return fmt.Errorf("accrual_system_address field: empty")
	}

	return nil
}

// NewDefaultConfig builds a Config with default values.
func NewDefaultConfig() Config {
	return Config{
		AccrualSysAddress: "http://127.0.0.1:8081",
	}
}
