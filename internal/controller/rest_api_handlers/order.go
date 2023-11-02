package restapihandlers

import (
	"context"
	"io"
	"net/http"

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
	var user models.User
	var orderNumberRaw []byte
	var ctx = context.Background()
	var err error

	orderNumberRaw, err = io.ReadAll(r.Body)
	if err != nil {
		errorJSON(w, err, or.logger)
		return
	}
	defer r.Body.Close()

	order.Number = string(orderNumberRaw)

	if !order.ValidLuhnNumber() {
		err = &models.AppError{Op: "handlersOrder.UploadOrders", Code: models.ErrOrderNumberLuhn}
		errorJSON(w, err, or.logger)
		return
	}

	user.Login, err = user.GetLoginFromContext(r.Context())
	if err != nil {
		errorJSON(w, err, or.logger)
		return
	}

	order.User = user.Login

	err = or.orderInteractor.Create(ctx, order)
	if err != nil {
		errorJSON(w, err, or.logger)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (or *orderRoutes) downloadOrders(w http.ResponseWriter, r *http.Request) {
	var login models.User
	var err error

	login.Login, err = login.GetLoginFromContext(r.Context())
	if err != nil {
		errorJSON(w, err, or.logger)
		return
	}

	orders, err := or.orderInteractor.GetList(r.Context(), login)
	if err != nil {
		errorJSON(w, err, or.logger)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	writeJSON(w, http.StatusOK, orders, or.logger)
}
