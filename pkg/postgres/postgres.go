package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Chystik/gophermart/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	connStr            = "host=%s port=%d user=%s password=%s sslmode=%s"
	connStrDB          = "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s"
	logDatabaseCreated = "database %s created"
)

type PgClient struct {
	db         *sqlx.DB
	connConfig *pgx.ConnConfig
	logger     logger.AppLogger
}

// opens a db and perform migrations
func NewPgClient(uri string, logger logger.AppLogger) (*PgClient, error) {
	cc, err := pgx.ParseURI(uri)
	if err != nil {
		return nil, err
	}

	if cc.Port == 0 {
		cc.Port = 5432
	}

	db, err := sqlx.Open("pgx", uri)
	if err != nil {
		return nil, err
	}

	return &PgClient{
		db:         db,
		connConfig: &cc,
		logger:     logger,
	}, nil
}

// Connect to a database and verify with a ping, if successful - create db if not exist
func (pc *PgClient) Connect(ctx context.Context) error {
	var err error
	var SSLmode string

	if pc.connConfig.TLSConfig == nil {
		SSLmode = "disable"
	}

	pc.db, err = sqlx.ConnectContext(
		ctx,
		"pgx",
		fmt.Sprintf(
			connStr,
			pc.connConfig.Host,
			pc.connConfig.Port,
			pc.connConfig.User,
			pc.connConfig.Password,
			SSLmode,
		),
	)
	if err != nil {
		pc.logger.Error(err.Error())
		return err
	}

	_, err = pc.db.Exec(fmt.Sprintf("CREATE DATABASE %s", pc.connConfig.Database))
	if err != nil {
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) || pgerrcode.DuplicateDatabase != pgErr.Code {
			pc.logger.Error(err.Error())
			return err
		}
		pc.logger.Info(err.Error())
	} else {
		pc.logger.Info(fmt.Sprintf(logDatabaseCreated, pc.connConfig.Database))
	}

	pc.db, err = sqlx.ConnectContext(
		ctx,
		"pgx",
		fmt.Sprintf(
			connStrDB,
			pc.connConfig.Host,
			pc.connConfig.Port,
			pc.connConfig.User,
			pc.connConfig.Password,
			pc.connConfig.Database,
			SSLmode,
		),
	)
	if err != nil {
		pc.logger.Error(err.Error())
		return err
	}

	return nil
}

func (pc *PgClient) Migrate() error {
	d, err := postgres.WithInstance(pc.db.DB, &postgres.Config{})
	if err != nil {
		pc.logger.Error(err.Error())
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://schema",
		pc.connConfig.Database, d)
	if err != nil {
		pc.logger.Error(err.Error())
		return err
	}

	err = m.Up()
	if err != nil && err.Error() != "no change" {
		pc.logger.Error(err.Error())
		return err
	}

	return nil
}

func (pc *PgClient) Disconnect(ctx context.Context) error {
	return pc.db.Close()
}

func (pc *PgClient) Ping(ctx context.Context) error {
	return pc.db.PingContext(ctx)
}

func (pc *PgClient) PingHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	w.Header().Set("Content-Tye", "text/plain")

	err := pc.Ping(ctx)
	if err != nil {
		pc.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (pc *PgClient) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return pc.db.ExecContext(ctx, query, args...)
}

func (pc *PgClient) GetContext(ctx context.Context, dest any, query string, args ...any) error {
	return pc.db.GetContext(ctx, dest, query, args...)
}

func (pc *PgClient) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return pc.db.QueryRowContext(ctx, query, args...)
}

func (pc *PgClient) SelectContext(ctx context.Context, dest any, query string, args ...any) error {
	return pc.db.SelectContext(ctx, dest, query, args...)
}
