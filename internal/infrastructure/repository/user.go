package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Chystik/gophermart/internal/models"
	"github.com/Chystik/gophermart/pkg/postgres"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNotFound = errors.New("not found in the repository")
	ErrExists   = errors.New("login already exists in the repository")
)

type userRepository struct {
	*postgres.PgClient
}

func NewUserRepository(db *postgres.PgClient) *userRepository {
	return &userRepository{db}
}

func (ur *userRepository) Create(ctx context.Context, user models.User) error {
	/* u, err := fromDomainUser(user)

	if err != nil {
		return err
	} */

	query := `
			INSERT INTO	praktikum.user (login, password, balance, withdrawn)
			VALUES ($1, $2, $3, $4)`

	_, err := ur.ExecContext(ctx, query, user.Login, user.Password, user.Balance, user.Withdrawn)
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
	var u models.User

	query := `
			SELECT login, password, balance, withdrawn
			FROM praktikum.user
			WHERE login = $1`

	err := ur.GetContext(ctx, &u, query, user.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, ErrNotFound
		} else {
			return u, err
		}
	}

	return u, nil
}
