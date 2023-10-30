package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Chystik/gophermart/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNotFound = errors.New("object not found")
	ErrExists   = errors.New("object already exists in the repository")
)

type userRepository struct {
	*sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *userRepository {
	return &userRepository{db}
}

func (ur *userRepository) Create(ctx context.Context, user models.User) error {
	query := `
			INSERT INTO	praktikum.user (login, password, balance, withdrawn)
			VALUES (:login, :password, :balance, :withdrawn)`

	_, err := sqlx.NamedExecContext(ctx, ur, query, user)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if !ok {
			return err
		} else if pgErr.Code == "23505" { // login exists: duplicate key value violates unique constraint
			return ErrExists
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

	rows, err := sqlx.NamedQueryContext(ctx, ur, query, user)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, ErrNotFound
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

	_, err := sqlx.NamedExecContext(ctx, ur, query, user)

	return err
}
