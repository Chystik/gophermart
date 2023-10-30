package usecase

import (
	"context"
	"time"

	"github.com/Chystik/gophermart/internal/infrastructure/repository"
	"github.com/Chystik/gophermart/internal/models"
)

type orderInteractor struct {
	orderRepo OrderRepository
	//accrual   AccrualWebAPI
}

func NewOrderInteractor(or OrderRepository) *orderInteractor {
	return &orderInteractor{
		orderRepo: or,
		//accrual:   a,
	}
}

func (oi *orderInteractor) Create(ctx context.Context, order models.Order) error {
	var err error

	o, err := oi.orderRepo.Get(ctx, order)
	if err != nil {
		if err == repository.ErrNotFound {
			order.Status = "NEW"
			order.UploadedAt = models.RFC3339Time{Time: time.Now()}
			return oi.orderRepo.Create(ctx, order)
		}
	}

	if o.User == order.User {
		return repository.ErrUploadedByUser
	}

	return repository.ErrUploadedByAnotherUser
}

func (oi *orderInteractor) GetList(ctx context.Context, login models.User) ([]models.Order, error) {
	return oi.orderRepo.GetList(ctx, login)
}
