package websocket

import (
	"fmt"

	"github.com/googollee/go-socket.io"
	"github.com/rs/zerolog/log"

	"github.com/vstdy/gophermart/service/gophermart"
)

// addRootNamespace sets actions and events for the root namespace.
func addRootNamespace(srv *socketio.Server, svc gophermart.Service, config Config) {
	srv.OnConnect(config.NotificationNamespace, onConnect)
	srv.OnError(config.NotificationNamespace, onError)
	srv.OnDisconnect(config.NotificationNamespace, onDisconnect)
	addRootEvents(srv, svc, config)
}

// onConnect performs actions for on the open event.
func onConnect(s socketio.Conn) error {
	go s.Emit("status", fmt.Sprintf("Connected! Session ID: %s", s.ID()))

	return nil
}

// onError performs actions on errors.
func onError(s socketio.Conn, err error) {
	log.Warn().Err(err).Msg("Websocket error:")
}

// onDisconnect performs actions on the disconnect event.
func onDisconnect(s socketio.Conn, reason string) {
	log.Info().Msgf("Websocket disconnect: %s", reason)
}
