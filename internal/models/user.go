package models

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Login     string `json:"login" db:"login"`
	Password  string `json:"password" db:"password"`
	Balance   Money  `db:"balance"`
	Withdrawn Money  `db:"withdrawn"`
}

type Withdrawal struct {
	Order       string      `json:"order" db:"order_id"`
	Sum         Money       `json:"sum" db:"sum"`
	ProcessedAt RFC3339Time `json:"processed_at" db:"processed_at"`
}

// SetPassword hashes the user's password
func (u *User) SetPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u User) Authenticate(actual User) error {
	err := bcrypt.CompareHashAndPassword([]byte(actual.Password), []byte(u.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return &AppError{Op: "user.Authenticate", Code: ErrUserCreds, Message: err.Error()}
		}
		return err
	}

	return nil
}

func (u User) GetLoginFromContext(ctx context.Context) (string, error) {
	claims, ok := ctx.Value(ClaimsKeyName).(*AuthClaims)
	if !ok {
		return "", &AppError{Op: "user.GetLoginFromContext", Code: ErrAuthClaims}
	}

	return claims.Login, nil
}
