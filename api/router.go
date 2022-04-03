package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"

	"github.com/vstdy0/go-diploma/cmd/gophermart/cmd/common"
	"github.com/vstdy0/go-diploma/service/gophermart"
)

// NewRouter returns router.
func NewRouter(svc gophermart.Service, config common.Config) chi.Router {
	h := NewHandler(svc, config)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Timeout(config.Timeout))
	r.Use(gzipDecompressRequest)
	r.Use(gzipCompressResponse)

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
				r.Post("/", h.addUserOrder)
				r.Get("/", h.getUserOrders)
			})

			r.Route("/balance", func(r chi.Router) {
				r.Get("/", h.getUserBalance)
				r.Post("/withdraw", h.addWithdrawal)
				r.Get("/withdrawals", h.getUserWithdrawals)
			})
		})
	})

	return r
}
