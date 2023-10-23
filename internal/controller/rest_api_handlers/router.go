package restapihandlers

import (
	md "github.com/Chystik/gophermart/internal/middleware"
	"github.com/Chystik/gophermart/internal/usecase"
	"github.com/Chystik/gophermart/pkg/logger"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(handler *chi.Mux, s usecase.GophermartInteractior, JWTkey []byte, l logger.AppLogger) {
	// Options
	handler.Use(md.MidLogger(l).WithLogging)
	//handler.Use(md.GzipMiddleware)
	handler.Use(middleware.Compress(5))
	handler.Use(middleware.Recoverer)

	// Routes
	//{
	newGophermartRoutes(handler, s, JWTkey, l)
	//}
}
