package usecase

import (
	"context"

	"github.com/Chystik/gophermart/internal/models"
)

type UserInteractor interface {
	Register(context.Context, models.User) error
	Authenticate(context.Context, models.User) error
	Get(context.Context, models.User) (models.User, error)
}

type OrderInteractor interface {
	Create(context.Context, models.Order) error
	GetAll(context.Context) ([]models.Order, error)
}

type UserRepository interface {
	Create(context.Context, models.User) error
	Get(context.Context, models.User) (models.User, error)
}

type OrderRepository interface {
	Create(context.Context, models.Order) error
	Get(context.Context, models.Order) (models.Order, error)
	GetList(context.Context) ([]models.Order, error)
}

type AccrualWebAPI interface {
	GetOrder(context.Context, models.Order) (models.Order, error)
}
