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

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

const (
	tokenExpiration = 5 * time.Minute
)

type gophermartRoutes struct {
	gophermartInteractor usecase.GophermartInteractior
	logger               logger.AppLogger
	JWTkey               []byte
}

func newGophermartRoutes(r *chi.Mux, s usecase.GophermartInteractior, JWTkey []byte, l logger.AppLogger) {
	h := &gophermartRoutes{
		gophermartInteractor: s,
		logger:               l,
		JWTkey:               JWTkey,
	}

	r.Route("/api/user/", func(r chi.Router) {
		// Public Routes
		r.Group(func(r chi.Router) {
			r.Post("/register", h.Register)
			r.Post("/login", h.Login)
		})

		// Private Routes
		// Require Authentication
		r.Group(func(r chi.Router) {
			//r.Use(AuthMiddleware)
			r.Route("/orders/", func(r chi.Router) {
				r.Post("/", h.UploadOrders)
				r.Get("/", h.DownloadOrders)
			})
			r.Route("/balance", func(r chi.Router) {
				r.Get("/", h.GetBalance)
				r.Post("/withdraw", h.Withdraw)
			})
			r.Get("/withdrawals", h.GetWithdrawals)
		})
	})

}

func (h *gophermartRoutes) Register(w http.ResponseWriter, r *http.Request) {
	//var creds credentials
	var user models.User
	var ctx = context.Background()
	var err error

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		errorJSON(w, err, http.StatusBadRequest, h.logger)
		return
	}

	// Hash user password
	user.SetPassword()

	// Create user
	err = h.gophermartInteractor.RegisterUser(ctx, user)
	if err != nil {
		// Login conflict
		if err == repository.ErrExists {
			errorJSON(w, err, http.StatusConflict, h.logger)
			return
		}
		// server err
		errorJSON(w, err, http.StatusInternalServerError, h.logger)
		return
	}

	// Authorize user
	h.authorize(w, user)
}

func (h *gophermartRoutes) Login(w http.ResponseWriter, r *http.Request) {
	//var creds credentials
	var user models.User
	var ctx = context.Background()
	var err error

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		errorJSON(w, err, http.StatusBadRequest, h.logger)
		return
	}

	// Verify password
	err = h.gophermartInteractor.AuthenticateUser(ctx, user)
	if err != nil {
		if err == repository.ErrNotFound { // Not found
			errorJSON(w, err, http.StatusUnauthorized, h.logger)
			return
		} else if err == models.ErrWrongCreds { // Wrong password
			errorJSON(w, err, http.StatusUnauthorized, h.logger)
			return
		}
		errorJSON(w, err, http.StatusInternalServerError, h.logger)
		return
	}

	// Authorize user
	h.authorize(w, user)
}

func (h *gophermartRoutes) UploadOrders(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	var ctx = context.Background()
	var err error

	err = json.NewDecoder(r.Body).Decode(&order.Number)
	if err != nil {
		errorJSON(w, err, http.StatusBadRequest, h.logger)
		return
	}

	claims, ok := r.Context().Value("props").(models.AuthClaims)
	if !ok {
		errorJSON(w, err, http.StatusUnauthorized, h.logger)
		return
	}
	order.User = claims.Login

	err = h.gophermartInteractor.CreateOrder(ctx, order)
	if err != nil {
		errorJSON(w, err, http.StatusBadRequest, h.logger)
		return
	}
}

func (h *gophermartRoutes) DownloadOrders(w http.ResponseWriter, r *http.Request) {

}

func (h *gophermartRoutes) GetBalance(w http.ResponseWriter, r *http.Request) {

}

func (h *gophermartRoutes) Withdraw(w http.ResponseWriter, r *http.Request) {

}

func (h *gophermartRoutes) GetWithdrawals(w http.ResponseWriter, r *http.Request) {

}

func (h *gophermartRoutes) authorize(w http.ResponseWriter, user models.User) {
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
	tokenStr, err := token.SignedString(h.JWTkey)
	if err != nil {
		errorJSON(w, err, http.StatusInternalServerError, h.logger)
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenStr,
		Expires: expirationTime,
	})
}
