package run

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Chystik/gophermart/config"
	restapihandlers "github.com/Chystik/gophermart/internal/controller/rest_api_handlers"
	"github.com/Chystik/gophermart/internal/infrastructure/repository"
	"github.com/Chystik/gophermart/internal/usecase"
	"github.com/Chystik/gophermart/pkg/httpserver"
	"github.com/Chystik/gophermart/pkg/logger"
	"github.com/Chystik/gophermart/pkg/postgres"

	"github.com/go-chi/chi/v5"
)

const (
	defaultLogLevel        = "info"
	defaultShutdownTimeout = 5 * time.Second

	logHTTPServerStart             = "HTTP server started on port: %s"
	logHTTPServerStop              = "Stopped serving new connections"
	logSignalInterrupt             = "Interrupt signal. Shutdown"
	logGracefulHTTPServerShutdown  = "Graceful shutdown of HTTP Server complete."
	logStorageSyncStart            = "data syncronization to file %s with interval %v started"
	logStorageSyncStop             = "Stopped saving storage data to a file"
	logGracefulStorageSyncShutdown = "Graceful shutdown of storage sync complete."
	logDBDisconnect                = "Graceful close connection for DB client complete."
)

func App(cfg *config.App, quit chan os.Signal) {

	// Logger
	logger, err := logger.Initialize(defaultLogLevel)
	if err != nil {
		panic(err)
	}

	// Postgres client
	pgClient, err := postgres.NewPgClient(cfg.DBuri, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	err = pgClient.Connect(ctx)
	if err != nil {
		logger.Fatal(err.Error())
	}

	err = pgClient.Migrate()
	if err != nil {
		logger.Fatal(err.Error())
	}

	// Repository
	userRepo := repository.NewUserRepository(pgClient)
	orderRepo := repository.NewOrderRepository(pgClient)

	// Interactor
	gophermartInteractor := usecase.NewGophermartInteractor(userRepo, orderRepo)

	// Router
	handler := chi.NewRouter()
	restapihandlers.NewRouter(handler, gophermartInteractor, cfg.JWTkey, logger)

	// HTTP client

	// HTTP server
	server := httpserver.NewServer(handler, httpserver.Address(cfg.Address))
	go func() {
		logger.Info(fmt.Sprintf(logHTTPServerStart, cfg.Address))
		if err := server.Startup(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err.Error())
		}
		logger.Info(logHTTPServerStop)
	}()

	// Wait signal
	<-quit
	logger.Info(logSignalInterrupt)
	ctxShutdown, shutdown := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer shutdown()

	// Graceful shutdown HTTP Server
	if err := server.Shutdown(ctxShutdown); err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info(logGracefulHTTPServerShutdown)

	// Graceful disconnect db client
	if err := pgClient.Disconnect(ctxShutdown); err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info(logDBDisconnect)
}
