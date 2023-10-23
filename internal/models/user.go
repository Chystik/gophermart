package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrWrongCreds = errors.New("mismatch login or password")
)

type User struct {
	Login     string  `json:"login" db:"login"`
	Password  string  `json:"password" db:"password"`
	Balance   float64 `db:"balance"`
	Withdrawn float64 `db:"withdrawn"`
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

func (u *User) Authenticate(actual User) error {
	err := bcrypt.CompareHashAndPassword([]byte(actual.Password), []byte(u.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return ErrWrongCreds
		}
		return err
	}

	return nil
}
