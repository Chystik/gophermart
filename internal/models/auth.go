package models

import "github.com/golang-jwt/jwt/v5"

type ClaimsKey string

type AuthClaims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}
