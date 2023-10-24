package usecase

import (
	"context"

	"github.com/Chystik/gophermart/internal/models"
)

type (
	UserInteractor interface {
		Register(context.Context, models.User) error
		Authenticate(context.Context, models.User) error
		Get(context.Context, models.User) (models.User, error)
	}

	OrderInteractor interface {
		Create(context.Context, models.Order) error
		GetAll(context.Context) ([]models.Order, error)
	}

	UserRepository interface {
		Create(context.Context, models.User) error
		Get(context.Context, models.User) (models.User, error)
	}

	OrderRepository interface {
		Create(context.Context, models.Order) error
		Get(context.Context, models.Order) (models.Order, error)
		GetList(context.Context) ([]models.Order, error)
	}

	AccrualWebAPI interface {
		GetOrder(context.Context, models.Order) (models.Order, error)
	}
)
