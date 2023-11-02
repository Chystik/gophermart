package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Chystik/gophermart/internal/models"

	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	*sqlx.DB
	getter *trmsqlx.CtxGetter
}

func NewUserRepository(db *sqlx.DB, getter *trmsqlx.CtxGetter) *userRepository {
	return &userRepository{
		DB:     db,
		getter: getter,
	}
}

func (ur *userRepository) Create(ctx context.Context, user models.User) error {
	query := `
			INSERT INTO	praktikum.user (login, password, balance, withdrawn)
			VALUES (:login, :password, :balance, :withdrawn)`

	_, err := sqlx.NamedExecContext(ctx, ur.getter.DefaultTrOrDB(ctx, ur.DB), query, user)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if !ok {
			return err
		} else if pgErr.Code == "23505" { // login exists: duplicate key value violates unique constraint
			return &models.AppError{Op: "userRepository.Create", Code: models.ErrExists, Message: "user already exists"}
		}
		return err
	}

	return nil
}

func (ur *userRepository) Get(ctx context.Context, user models.User) (models.User, error) {
	query := `
			SELECT login, password, balance, withdrawn
			FROM praktikum.user
			WHERE login = :login`

	rows, err := sqlx.NamedQueryContext(ctx, ur.getter.DefaultTrOrDB(ctx, ur.DB), query, user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, &models.AppError{Op: "userRepository.Get", Code: models.ErrNotFound, Message: "user not found"}
		} else {
			return user, err
		}
	}
	defer rows.Close()

	// Expect only one result
	if !rows.Next() {
		return user, rows.Err()
	}

	err = rows.StructScan(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (ur *userRepository) Update(ctx context.Context, user models.User) error {
	query := `
			UPDATE praktikum.user 
			SET balance = :balance, withdrawn = :withdrawn
			WHERE login = :login`

	_, err := sqlx.NamedExecContext(ctx, ur.getter.DefaultTrOrDB(ctx, ur.DB), query, user)

	return err
}
