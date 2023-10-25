package restapihandlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/Chystik/gophermart/internal/infrastructure/repository"
	"github.com/Chystik/gophermart/internal/models"
	"github.com/Chystik/gophermart/internal/usecase"
	"github.com/Chystik/gophermart/pkg/logger"
)

var (
	errNotValidLuhn = errors.New("not valid order number")
)

type orderRoutes struct {
	orderInteractor usecase.OrderInteractor
	logger          logger.AppLogger
}

func newOrderRoutes(oi usecase.OrderInteractor, l logger.AppLogger) *orderRoutes {
	return &orderRoutes{
		orderInteractor: oi,
		logger:          l,
	}
}

func (or *orderRoutes) uploadOrders(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	var numberRaw []byte
	var number int
	var ctx = context.Background()
	var err error

	numberRaw, err = io.ReadAll(r.Body)
	if err != nil {
		errorJSON(w, err, http.StatusInternalServerError, or.logger)
		return
	}

	// check number
	number, err = strconv.Atoi(string(numberRaw))
	if err != nil {
		errorJSON(w, err, http.StatusUnprocessableEntity, or.logger)
		return
	}

	// validate
	if !valid(number) {
		errorJSON(w, errNotValidLuhn, http.StatusUnprocessableEntity, or.logger)
		return
	}

	order.Number = string(numberRaw)
	order.User, err = getUserLogin(r.Context()) //claims.Login
	if err != nil {
		errorJSON(w, err, http.StatusUnauthorized, or.logger)
		return
	}

	err = or.orderInteractor.Create(ctx, order)
	if err != nil {
		if err == repository.ErrUploadedByUser {
			w.WriteHeader(http.StatusOK)
			return
		} else if err == repository.ErrUploadedByAnotherUser {
			w.WriteHeader(http.StatusConflict)
			return
		}
		errorJSON(w, err, http.StatusInternalServerError, or.logger)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (or *orderRoutes) downloadOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := or.orderInteractor.GetAll(r.Context())
	if err != nil {
		errorJSON(w, err, http.StatusInternalServerError, or.logger)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	writeJSON(w, http.StatusOK, "application/json", orders, or.logger)
}

// Valid checks the order number using the Luhn algorithm
func valid(number int) bool {
	var luhn int

	checksum := func(n int) int {
		for i := 0; n > 0; i++ {
			cur := n % 10

			if i%2 == 0 {
				cur = cur * 2
				if cur > 9 {
					cur = cur%10 + cur/10
				}
			}

			luhn += cur
			n = n / 10
		}
		return luhn % 10
	}

	return (number%10+checksum(number/10))%10 == 0
}
