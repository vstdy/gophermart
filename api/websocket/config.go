package websocket

import (
	"fmt"

	"github.com/go-chi/jwtauth/v5"
)

// Config keeps websocket params.
type Config struct {
	NotificationNamespace string           `mapstructure:"notification_namespace"`
	JWTAuth               *jwtauth.JWTAuth `mapstructure:"-"`
}

// Validate performs a basic validation.
func (config Config) Validate() error {
	if config.NotificationNamespace == "" {
		return fmt.Errorf("notification_namespace field: empty")
	}

	return nil
}

// NewDefaultConfig builds a Config with default values.
func NewDefaultConfig() Config {
	return Config{
		NotificationNamespace: "/",
	}
}
