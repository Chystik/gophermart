package restapihandlers

import (
	md "github.com/Chystik/gophermart/internal/middleware"
	"github.com/Chystik/gophermart/internal/usecase"
	"github.com/Chystik/gophermart/pkg/logger"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(handler *chi.Mux, ui usecase.UserInteractor, oi usecase.OrderInteractor, JWTkey []byte, l logger.AppLogger) {
	// Options
	handler.Use(md.MidLogger(l).WithLogging)
	handler.Use(md.GzipMiddleware)
	//handler.Use(middleware.Compress(5))
	handler.Use(middleware.Recoverer)

	// Routes
	ur := newUserRoutes(ui, JWTkey, l)
	or := newOrderRoutes(oi, l)

	handler.Route("/api/user", func(r chi.Router) {
		// Public Routes
		r.Group(func(r chi.Router) {
			r.Post("/register", ur.register)
			r.Post("/login", ur.login)
		})

		// Private Routes
		// Require Authentication
		r.Group(func(r chi.Router) {
			r.Use(md.Authentication)
			r.Route("/orders", func(r chi.Router) {
				r.Post("/", or.uploadOrders)
				r.Get("/", or.downloadOrders)
			})
			r.Route("/balance", func(r chi.Router) {
				r.Get("/", ur.getBalance)
				r.Post("/withdraw", ur.withdraw)
			})
			r.Get("/withdrawals", ur.getWithdrawals)
		})
	})
}
