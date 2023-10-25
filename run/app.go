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
	"github.com/Chystik/gophermart/internal/infrastructure/webapi"
	"github.com/Chystik/gophermart/internal/usecase"
	"github.com/Chystik/gophermart/internal/usecase/syncer"
	"github.com/Chystik/gophermart/pkg/httpclient"
	"github.com/Chystik/gophermart/pkg/httpserver"
	"github.com/Chystik/gophermart/pkg/logger"
	"github.com/Chystik/gophermart/pkg/postgres"

	"github.com/go-chi/chi/v5"
)

const (
	defaultLogLevel        = "info"
	defaultShutdownTimeout = 5 * time.Second
	accrualAddrScheme      = "http"

	logHTTPServerStart            = "HTTP server started on port: %s"
	logHTTPServerStop             = "Stopped serving new connections"
	logSignalInterrupt            = "Interrupt signal. Shutdown"
	logGracefulHTTPServerShutdown = "Graceful shutdown of HTTP Server complete."
	logSyncStart                  = "data syncronization with accrual service %s started"
	logSyncStop                   = "Stopped syncronization with accrual service"
	logGracefulSyncShutdown       = "Graceful shutdown of accrual sync complete."
	logDBDisconnect               = "Graceful close connection for DB client complete."
)

func App(cfg *config.App, quit chan os.Signal) {
	// Logger
	logger, err := logger.Initialize(defaultLogLevel, "app.log")
	if err != nil {
		panic(err)
	}

	logger.Info(fmt.Sprintf("%#v", cfg))

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

	// HTTP client
	httpClient := httpclient.NewClient(httpclient.Timeout(20 * time.Second))

	// Accrual web API
	accrualWebAPI := webapi.NewAccrualWebAPI(httpClient, webapi.Address(cfg.AccrualAddress))

	// Repository
	userRepo := repository.NewUserRepository(pgClient)
	orderRepo := repository.NewOrderRepository(pgClient)
	withdrawalRepo := repository.NewWithdrawalRepository(pgClient)

	// Interactor
	userInteractor := usecase.NewUserInteractor(userRepo, withdrawalRepo)
	orderInteractor := usecase.NewOrderInteractor(orderRepo)

	// Syncer
	syncer := syncer.NewSyncer(userRepo, orderRepo, accrualWebAPI, logger)

	// Router
	handler := chi.NewRouter()
	restapihandlers.NewRouter(handler, userInteractor, orderInteractor, cfg.JWTkey, logger)

	// HTTP server
	server := httpserver.NewServer(handler, httpserver.Address(cfg.Address))
	go func() {
		logger.Info(fmt.Sprintf(logHTTPServerStart, cfg.Address))
		if err := server.Startup(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err.Error())
		}
		logger.Info(logHTTPServerStop)
	}()

	// Start syncer
	go func() {
		logger.Info(fmt.Sprintf(logSyncStart, cfg.AccrualAddress))
		if err := syncer.Run(); err != nil {
			logger.Fatal(err.Error())
		}
		logger.Info(logSyncStop)
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

	// Graceful shutdown syncer
	if err := syncer.Shutdown(ctxShutdown); err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info(logGracefulSyncShutdown)
}
