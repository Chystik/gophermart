package usecase

import (
	"context"

	"github.com/Chystik/gophermart/internal/infrastructure/repository"
	"github.com/Chystik/gophermart/internal/models"
)

type orderInteractor struct {
	orderRepo OrderRepository
	accrual   AccrualWebAPI
}

func NewOrderInteractor(or OrderRepository, a AccrualWebAPI) *orderInteractor {
	return &orderInteractor{
		orderRepo: or,
		accrual:   a,
	}
}

func (oi *orderInteractor) Create(ctx context.Context, order models.Order) error {
	var err error

	// call webapi
	order, err = oi.accrual.GetOrder(ctx, order)
	if err != nil {
		// TODO check status
		return err
	}

	o, err := oi.orderRepo.Get(ctx, order)
	if err != nil {
		if err == repository.ErrNotFound {
			return oi.orderRepo.Create(ctx, order)
		}
	}

	if o.User == order.User {
		return repository.ErrUploadedByUser
	}

	return repository.ErrUploadedByAnotherUser
}

func (oi *orderInteractor) GetAll(ctx context.Context) ([]models.Order, error) {
	return oi.orderRepo.GetList(ctx)
}
