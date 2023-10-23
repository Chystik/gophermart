package usecase

import (
	"context"
	"time"

	"github.com/Chystik/gophermart/internal/infrastructure/repository"
	"github.com/Chystik/gophermart/internal/models"
)

type orderInteractor struct {
	orderRepo OrderRepository
}

func NewOrderInteractor(or OrderRepository) *orderInteractor {
	return &orderInteractor{
		orderRepo: or,
	}
}

func (oi *orderInteractor) Create(ctx context.Context, order models.Order) error {
	// call webapi
	order.Status = "NEW"
	order.UploadedAt = time.Now()

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
