package restapihandlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Chystik/gophermart/internal/infrastructure/repository"
	"github.com/Chystik/gophermart/internal/models"
	"github.com/Chystik/gophermart/internal/usecase"
	"github.com/Chystik/gophermart/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
)

const (
	cookieName                       = "token"
	claimsKey       models.ClaimsKey = "props"
	tokenExpiration                  = 5 * time.Minute
)

var (
	errWrongAuthClaims = errors.New("wrong auth claims")
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
	err = user.SetPassword()
	if err != nil {
		errorJSON(w, err, http.StatusInternalServerError, ur.logger)
		return
	}

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
	w.WriteHeader(http.StatusOK)
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
	w.WriteHeader(http.StatusOK)
}

func (ur *userRoutes) getBalance(w http.ResponseWriter, r *http.Request) {
	var ctx = context.Background()

	login, err := getUserLogin(r.Context())
	if err != nil {
		errorJSON(w, err, http.StatusUnauthorized, ur.logger)
		return
	}

	user := models.User{
		Login: login,
	}

	user, err = ur.userInteractor.Get(ctx, user)
	if err != nil {
		errorJSON(w, err, http.StatusInternalServerError, ur.logger)
		return
	}

	balance := fromDomainBalance(user)

	writeJSON(w, http.StatusOK, "application/json", balance, ur.logger)
}

func (ur *userRoutes) withdraw(w http.ResponseWriter, r *http.Request) {
	var wth models.Withdrawal
	var ctx = context.Background()
	var order int
	var err error
	var login string

	err = json.NewDecoder(r.Body).Decode(&wth)
	if err != nil {
		errorJSON(w, err, http.StatusUnprocessableEntity, ur.logger)
		return
	}

	order, err = strconv.Atoi(wth.Order)
	if err != nil {
		errorJSON(w, errNotValidLuhn, http.StatusUnprocessableEntity, ur.logger)
		return
	}

	// Validate order number
	if !valid(order) {
		errorJSON(w, errNotValidLuhn, http.StatusUnprocessableEntity, ur.logger)
		return
	}

	login, err = getUserLogin(r.Context())
	if err != nil {
		errorJSON(w, errNotValidLuhn, http.StatusUnauthorized, ur.logger)
		return
	}

	err = ur.userInteractor.Withdraw(ctx, wth, models.User{Login: login})
	if err != nil {
		var stat int
		if err == usecase.ErrNotEnoughMoney {
			stat = http.StatusPaymentRequired
		} else {
			stat = http.StatusInternalServerError
		}
		errorJSON(w, err, stat, ur.logger)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ur *userRoutes) getWithdrawals(w http.ResponseWriter, r *http.Request) {
	var ctx = context.Background()

	wth, err := ur.userInteractor.GetWithdrawals(ctx)
	if err != nil {
		errorJSON(w, err, http.StatusInternalServerError, ur.logger)
		return
	}

	if len(wth) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	writeJSON(w, http.StatusOK, "application/json", wth, ur.logger)
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
		Name:    cookieName,
		Value:   tokenStr,
		Expires: expirationTime,
	})
}

func getUserLogin(ctx context.Context) (string, error) {
	claims, ok := ctx.Value(claimsKey).(*models.AuthClaims)
	if !ok {
		return "", errWrongAuthClaims
	}

	return claims.Login, nil
}
