package usecase

import (
	"context"
	"time"

	"github.com/Chystik/gophermart/internal/infrastructure/repository"
	"github.com/Chystik/gophermart/internal/models"
	"github.com/avito-tech/go-transaction-manager/trm"
)

type orderInteractor struct {
	orderRepo OrderRepository
	trm       trm.Manager
}

func NewOrderInteractor(or OrderRepository, trm trm.Manager) *orderInteractor {
	return &orderInteractor{
		orderRepo: or,
		trm:       trm,
	}
}

func (oi *orderInteractor) Create(ctx context.Context, order models.Order) error {
	var err error
	var o models.Order

	err = oi.trm.Do(ctx, func(ctx context.Context) error {
		o, err = oi.orderRepo.Get(ctx, order)
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
	})

	return err
}

func (oi *orderInteractor) GetList(ctx context.Context, login models.User) ([]models.Order, error) {
	var o []models.Order
	var err error

	err = oi.trm.Do(ctx, func(ctx context.Context) error {
		o, err = oi.orderRepo.GetList(ctx, login)
		return err
	})

	return o, err
}
