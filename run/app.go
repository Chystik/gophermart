package run

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/Chystik/gophermart/config"
	handlers "github.com/Chystik/gophermart/internal/controller/rest_api_handlers"
	"github.com/Chystik/gophermart/internal/infrastructure/repository"
	"github.com/Chystik/gophermart/internal/infrastructure/webapi"
	"github.com/Chystik/gophermart/internal/usecase"
	"github.com/Chystik/gophermart/internal/usecase/syncer"
	"github.com/Chystik/gophermart/pkg/httpclient"
	"github.com/Chystik/gophermart/pkg/httpserver"
	"github.com/Chystik/gophermart/pkg/logger"
	"github.com/Chystik/gophermart/pkg/postgres"

	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const (
	defaultLogLevel        = "info"
	defaultShutdownTimeout = 5 * time.Second
	accrualAddrScheme      = "http"

	logHTTPServerStart            = "HTTP server started"
	logHTTPServerStop             = "Stopped serving new connections"
	logSignalInterrupt            = "Interrupt signal. Shutdown"
	logGracefulHTTPServerShutdown = "Graceful shutdown of HTTP Server complete."
	logSyncStart                  = "data syncronization with accrual service started"
	logSyncStop                   = "Stopped syncronization with accrual service"
	logGracefulSyncShutdown       = "Graceful shutdown of accrual sync complete."
	logDBDisconnect               = "Graceful close connection for DB client complete."
)

func App(cfg *config.App, quit chan os.Signal) {
	// Logger
	logger, err := logger.Initialize(defaultLogLevel, "./app.log")
	if err != nil {
		panic(err)
	}

	// Postgres db
	pg, err := postgres.New(cfg.DBuri.String(), logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	err = pg.Connect(ctx)
	if err != nil {
		logger.Fatal(err.Error())
	}

	err = pg.Migrate()
	if err != nil {
		logger.Fatal(err.Error())
	}

	// HTTP client
	httpClient := httpclient.NewClient(httpclient.Timeout(20 * time.Second))

	// Accrual web API
	accrualWebAPI := webapi.NewAccrualWebAPI(httpClient, webapi.Address(cfg.AccrualAddress.String()))

	// Repository
	trManager := manager.Must(trmsqlx.NewDefaultFactory(pg.DB))
	userRepo := repository.NewUserRepository(pg.DB, trmsqlx.DefaultCtxGetter)
	orderRepo := repository.NewOrderRepository(pg.DB, trmsqlx.DefaultCtxGetter)
	withdrawalRepo := repository.NewWithdrawalRepository(pg.DB, trmsqlx.DefaultCtxGetter)

	// Interactor
	userInteractor := usecase.NewUserInteractor(userRepo, withdrawalRepo, trManager)
	orderInteractor := usecase.NewOrderInteractor(orderRepo, trManager)

	// Syncer
	syncer := syncer.NewSyncer(userRepo, orderRepo, accrualWebAPI, logger)

	// Router
	handler := chi.NewRouter()
	handlers.NewRouter(handler, userInteractor, orderInteractor, cfg.JWTkey, logger)

	// HTTP server
	server := httpserver.NewServer(handler, httpserver.Address(cfg.Address.String()))
	go func() {
		logger.Info(logHTTPServerStart, zap.String("addr", string(cfg.Address)))
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err.Error())
		}
		logger.Info(logHTTPServerStop)
	}()

	// Start syncer
	go func() {
		logger.Info(logSyncStart, zap.String("addr", string(cfg.AccrualAddress)))
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
	if err := pg.Disconnect(ctxShutdown); err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info(logDBDisconnect)

	// Graceful shutdown syncer
	if err := syncer.Shutdown(ctxShutdown); err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info(logGracefulSyncShutdown)
}
