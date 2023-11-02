package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ClaimsKey string

const (
	CookieName                = "token"
	ClaimsKeyName   ClaimsKey = "props"
	TokenExpiration           = 5 * time.Minute
)

type AuthClaims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}
