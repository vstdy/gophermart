package psql

import "fmt"

const (
	defaultConfigEndpoint = "postgres://user:password@localhost:5432/gophermart?sslmode=disable"
)

// Config keeps Storage configuration.
type Config struct {
	URI string `mapstructure:"database_uri"`
}

// Validate performs a basic validation.
func (config Config) Validate() error {
	if config.URI == "" {
		return fmt.Errorf("%s field: empty", "DSN")
	}

	return nil
}

// NewDefaultConfig builds a Config with default values.
func NewDefaultConfig() Config {
	return Config{
		URI: defaultConfigEndpoint,
	}
}
