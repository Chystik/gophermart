package restapihandlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Chystik/gophermart/internal/infrastructure/repository"
	"github.com/Chystik/gophermart/internal/models"
	"github.com/Chystik/gophermart/internal/usecase"
	"github.com/Chystik/gophermart/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
)

const (
	key             models.ClaimsKey = "props"
	tokenExpiration                  = 5 * time.Minute
)

type userRoutes struct {
	userInteractor usecase.UserInteractor
	logger         logger.AppLogger
	JWTkey         []byte
}

func newUserRoutes(ui usecase.UserInteractor, JWTkey []byte, l logger.AppLogger) *userRoutes {
	return &userRoutes{
		userInteractor: ui,
		logger:         l,
		JWTkey:         JWTkey,
	}
}

func (ur *userRoutes) register(w http.ResponseWriter, r *http.Request) {
	//var creds credentials
	var user models.User
	var ctx = context.Background()
	var err error

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		errorJSON(w, err, http.StatusBadRequest, ur.logger)
		return
	}

	// Hash user password
	user.SetPassword()

	// Create user
	err = ur.userInteractor.Register(ctx, user)
	if err != nil {
		// Login conflict
		if err == repository.ErrExists {
			errorJSON(w, err, http.StatusConflict, ur.logger)
			return
		}
		// server err
		errorJSON(w, err, http.StatusInternalServerError, ur.logger)
		return
	}

	// Authorize user
	ur.authorize(w, user)
}

func (ur *userRoutes) login(w http.ResponseWriter, r *http.Request) {
	//var creds credentials
	var user models.User
	var ctx = context.Background()
	var err error

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		errorJSON(w, err, http.StatusBadRequest, ur.logger)
		return
	}

	// Verify password
	err = ur.userInteractor.Authenticate(ctx, user)
	if err != nil {
		if err == repository.ErrNotFound { // Not found
			errorJSON(w, err, http.StatusUnauthorized, ur.logger)
			return
		} else if err == models.ErrWrongCreds { // Wrong password
			errorJSON(w, err, http.StatusUnauthorized, ur.logger)
			return
		}
		errorJSON(w, err, http.StatusInternalServerError, ur.logger)
		return
	}

	// Authorize user
	ur.authorize(w, user)
}

func (ur *userRoutes) getBalance(w http.ResponseWriter, r *http.Request) {

}

func (ur *userRoutes) withdraw(w http.ResponseWriter, r *http.Request) {

}

func (ur *userRoutes) getWithdrawals(w http.ResponseWriter, r *http.Request) {

}

func (ur *userRoutes) authorize(w http.ResponseWriter, user models.User) {
	expirationTime := time.Now().Add(tokenExpiration)

	// Create the JWT claims, which includes the username and expiry time
	claims := &models.AuthClaims{
		Login: user.Login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenStr, err := token.SignedString(ur.JWTkey)
	if err != nil {
		errorJSON(w, err, http.StatusInternalServerError, ur.logger)
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenStr,
		Expires: expirationTime,
	})
}
