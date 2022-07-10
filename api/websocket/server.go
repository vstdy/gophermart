package websocket

import (
	"github.com/googollee/go-socket.io"

	"github.com/vstdy/gophermart/service/gophermart"
)

// NewServer creates a new WebSocket-server.
func NewServer(svc gophermart.Service, config Config) *socketio.Server {
	server := socketio.NewServer(nil)
	addRootNamespace(server, svc, config)

	return server
}
