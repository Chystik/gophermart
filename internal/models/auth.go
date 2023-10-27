package models

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ClaimsKey string

const (
	CookieName                = "token"
	ClaimsKeyName   ClaimsKey = "props"
	TokenExpiration           = 5 * time.Minute
)

var (
	errWrongAuthClaims = errors.New("wrong auth claims")
)

type AuthClaims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}
