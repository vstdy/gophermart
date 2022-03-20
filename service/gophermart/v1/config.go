package gophermart

import (
	"fmt"
	"time"
)

// Config keeps Service params.
type Config struct {
	AccrualSysAddress   string        `mapstructure:"accrual_system_address"`
	UpdaterTimeout      time.Duration `mapstructure:"updater_timeout"`
	StatusCheckInterval time.Duration `mapstructure:"status_check_interval"`
}

// Validate performs a basic validation.
func (config Config) Validate() error {
	if config.UpdaterTimeout < time.Second {
		return fmt.Errorf("updater_timeout field: too short period")
	}

	if config.StatusCheckInterval < time.Second {
		return fmt.Errorf("status_check_interval field: too short period")
	}

	return nil
}

// NewDefaultConfig builds a Config with default values.
func NewDefaultConfig() Config {
	return Config{
		AccrualSysAddress:   "http://127.0.0.1:8081",
		UpdaterTimeout:      5 * time.Second,
		StatusCheckInterval: 5 * time.Second,
	}
}
