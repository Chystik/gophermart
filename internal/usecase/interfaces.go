package usecase

import (
	"context"

	"github.com/Chystik/gophermart/internal/models"
)

type GophermartInteractior interface {
	RegisterUser(context.Context, models.User) error
	GetUser(context.Context, models.User) (models.User, error)
	AuthenticateUser(context.Context, models.User) error
	CreateOrder(context.Context, models.Order) error
	GetOrders(context.Context) ([]models.Order, error)
}

type OrderRepository interface {
	Create(context.Context, models.Order) error
	GetList(context.Context) ([]models.Order, error)
}

type UserRepository interface {
	Create(context.Context, models.User) error
	Get(context.Context, models.User) (models.User, error)
}
