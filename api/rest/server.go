package rest

import (
	"net/http"

	"github.com/googollee/go-socket.io"

	"github.com/vstdy/gophermart/service/gophermart"
)

// NewServer creates a new HTTP-server.
func NewServer(svc gophermart.Service, ws *socketio.Server, config Config) *http.Server {
	router := NewRouter(svc, ws, config)

	return &http.Server{Addr: config.RunAddress, Handler: router}
}
