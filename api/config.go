package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/googollee/go-socket.io"
	"github.com/lestrrat-go/jwx/jwa"

	"github.com/vstdy/gophermart/api/rest"
	"github.com/vstdy/gophermart/api/websocket"
	"github.com/vstdy/gophermart/service/gophermart"
)

// Config keeps api params.
type Config struct {
	SecretKey string           `mapstructure:"secret_key"`
	Rest      rest.Config      `mapstructure:"rest,squash"`
	Websocket websocket.Config `mapstructure:"websocket,squash"`
}

// BuildAPI builds API dependency.
func (config Config) BuildAPI(svc gophermart.Service) (*http.Server, *socketio.Server) {
	config.fillCommonConfigs()
	ws := websocket.NewServer(svc, config.Websocket)
	r := rest.NewServer(svc, ws, config.Rest)

	return r, ws
}

// Validate performs a basic validation.
func (config Config) Validate() error {
	if config.SecretKey == "" {
		return fmt.Errorf("secret_key field: empty")
	}

	if err := config.Rest.Validate(); err != nil {
		return fmt.Errorf("provider: %w", err)
	}

	if err := config.Websocket.Validate(); err != nil {
		return fmt.Errorf("provider: %w", err)
	}

	return nil
}

// NewDefaultConfig builds a Config with default values.
func NewDefaultConfig() Config {
	return Config{
		SecretKey: "secret_key",
		Rest:      rest.NewDefaultConfig(),
		Websocket: websocket.NewDefaultConfig(),
	}
}

// fillCommonConfigs fills in common configs for http and websocket servers.
func (config *Config) fillCommonConfigs() {
	jwtAuth := jwtauth.New(jwa.HS256.String(), []byte(config.SecretKey), nil)
	config.Rest.JWTAuth = jwtAuth
	config.Websocket.JWTAuth = jwtAuth
}
