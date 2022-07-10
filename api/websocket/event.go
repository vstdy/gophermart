package websocket

import (
	"context"
	"fmt"

	"github.com/go-chi/jwtauth/v5"
	"github.com/googollee/go-socket.io"
	"github.com/rs/zerolog/log"

	"github.com/vstdy/gophermart/service/gophermart"
)

// addRootEvents adds events for the root namespace.
func addRootEvents(srv *socketio.Server, svc gophermart.Service, config Config) {
	srv.OnEvent(config.NotificationNamespace, "join_room", joinRoomEvent(config))
	srv.OnEvent(config.NotificationNamespace, "leave_room", leaveRoomEvent)
	go notifyUser(srv, svc, config)
}

// joinRoomEvent connects the user to the room.
func joinRoomEvent(config Config) func(s socketio.Conn, msg string) {
	return func(s socketio.Conn, token string) {
		jwt, err := verifyToken(config.JWTAuth, token)
		if err != nil {
			s.Emit("error", fmt.Sprintf("sign in: %s", err))
			return
		}

		ctx, ok := s.Context().(context.Context)
		if ok {
			userID := ctx.Value(jwtauth.TokenCtxKey)
			s.Leave(userID.(string))
		}

		userID, err := getUserID(jwt)
		if err != nil {
			s.Emit("error", fmt.Sprintf("sign in: %s", err))
			return
		}

		ctx = context.WithValue(context.Background(), jwtauth.TokenCtxKey, userID)
		s.SetContext(ctx)
		s.Join(userID)
		s.Emit("status", "Signed in!")
	}
}

// leaveRoomEvent disconnects the user from the room.
func leaveRoomEvent(s socketio.Conn) {
	ctx, ok := s.Context().(context.Context)
	if !ok {
		s.Emit("error", "sign out: unauthorized")
		return
	}
	userID := ctx.Value(jwtauth.TokenCtxKey)
	s.Leave(userID.(string))
	s.Emit("status", "Signed out")
}

// notifyUser sends a notification to the user's room.
func notifyUser(srv *socketio.Server, svc gophermart.Service, config Config) {
	ch := svc.GetAccrualNotificationsChan()
	for transaction := range ch {
		msg := fmt.Sprintf("You get %.2f bonuses for order %s", transaction.Accrual, transaction.Order)
		srv.BroadcastToRoom(config.NotificationNamespace, transaction.UserID.String(), "new_notification", msg)
	}

	log.Info().Msg("AccrualNotificationsChan closed")
}
