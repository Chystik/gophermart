package middleware

import (
	"context"
	"net/http"

	"github.com/Chystik/gophermart/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

var (
	JWTkey                      = []byte("my_secret_key")
	key        models.ClaimsKey = "props"
	cookieName                  = "token"
)

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(cookieName)
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tokenStr := c.Value
		claims := &models.AuthClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			return JWTkey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if token.Valid {
			ctx := context.WithValue(r.Context(), key, claims)
			// Access context values in handlers like this
			// props, _ := r.Context().Value("props").(*jwt.MapClaims)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	})
}
