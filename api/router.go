package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"

	"github.com/vstdy/gophermart/cmd/gophermart/cmd/common"
	"github.com/vstdy/gophermart/service/gophermart"
)

// NewRouter returns router.
func NewRouter(svc gophermart.Service, config common.Config) chi.Router {
	h := NewHandler(svc, config.SecretKey)
	r := chi.NewRouter()

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

	r.Route("/api/user", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.AllowContentType("application/json"))

			r.Post("/register", h.register)
			r.Post("/login", h.login)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(h.tokenAuth))
			r.Use(jwtauth.Authenticator)

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

	return r
}
