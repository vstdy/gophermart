package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	socketio "github.com/googollee/go-socket.io"

	"github.com/vstdy/gophermart/service/gophermart"
)

// NewRouter returns router.
func NewRouter(svc gophermart.Service, ws *socketio.Server, config Config) chi.Router {
	h := NewHandler(svc, config.JWTAuth)
	r := chi.NewRouter()

	r.Route("/api/user", func(r chi.Router) {
		r.Use(
			middleware.RequestID,
			middleware.RealIP,
			middleware.Logger,
			middleware.Recoverer,
			middleware.StripSlashes,
			middleware.Timeout(config.Timeout),
			gzipDecompressRequest,
			gzipCompressResponse,
		)

		// Public routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.AllowContentType("application/json"))

			r.Post("/register", h.register)
			r.Post("/login", h.login)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(
				jwtauth.Verifier(config.JWTAuth),
				jwtauth.Authenticator,
			)

			r.Route("/orders", func(r chi.Router) {
				r.Post("/", h.addUsersOrder)
				r.Get("/", h.getUsersOrders)
			})

			r.Route("/balance", func(r chi.Router) {
				r.Get("/", h.getUsersBalance)
				r.Post("/withdraw", h.addWithdrawal)
				r.Get("/withdrawals", h.getUsersWithdrawals)
			})
		})
	})

	r.Route("/notifications/socket", func(r chi.Router) {
		r.Handle("/", ws)
		r.Method(
			http.MethodGet, "/page",
			http.StripPrefix(
				"/notifications/socket/page",
				http.FileServer(http.Dir("./templates")),
			),
		)
	})

	return r
}
