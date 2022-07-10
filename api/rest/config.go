package rest

import (
	"fmt"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

// Config keeps rest params.
type Config struct {
	Timeout    time.Duration    `mapstructure:"timeout"`
	JWTAuth    *jwtauth.JWTAuth `mapstructure:"-"`
	RunAddress string           `mapstructure:"run_address"`
}

// Validate performs a basic validation.
func (config Config) Validate() error {
	if config.Timeout < time.Second {
		return fmt.Errorf("timeout field: must be GTE 1")
	}

	if config.RunAddress == "" {
		return fmt.Errorf("run_address field: empty")
	}

	return nil
}

// NewDefaultConfig builds a Config with default values.
func NewDefaultConfig() Config {
	return Config{
		Timeout:    5 * time.Second,
		RunAddress: "0.0.0.0:8080",
	}
}
