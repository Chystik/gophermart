package restapihandlers

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/Chystik/gophermart/internal/infrastructure/repository"
	"github.com/Chystik/gophermart/internal/models"
	"github.com/Chystik/gophermart/internal/usecase"
	"github.com/Chystik/gophermart/pkg/logger"
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
	var ctx = context.Background()
	var err error

	numberRaw, err := io.ReadAll(r.Body)
	if err != nil {
		errorJSON(w, err, http.StatusInternalServerError, or.logger)
		return
	}

	// check number
	_, err = strconv.Atoi(string(numberRaw))
	if err != nil {
		errorJSON(w, err, http.StatusUnprocessableEntity, or.logger)
		return
	}

	order.Number = string(numberRaw)

	claims, ok := r.Context().Value(key).(*models.AuthClaims)
	if !ok {
		errorJSON(w, err, http.StatusUnauthorized, or.logger)
		return
	}

	order.User = claims.Login

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
	var ctx = context.Background()

	orders, err := or.orderInteractor.GetAll(ctx)
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
