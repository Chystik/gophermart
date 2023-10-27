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
		Withdraw(context.Context, models.Withdrawal, models.User) error
		GetWithdrawals(context.Context) ([]models.Withdrawal, error)
	}

	OrderInteractor interface {
		Create(context.Context, models.Order) error
		GetList(context.Context, models.User) ([]models.Order, error)
	}

	UserRepository interface {
		Create(context.Context, models.User) error
		Get(context.Context, models.User) (models.User, error)
		Update(context.Context, models.User) error
	}

	OrderRepository interface {
		Create(context.Context, models.Order) error
		Get(context.Context, models.Order) (models.Order, error)
		GetList(context.Context, models.User) ([]models.Order, error)
		GetUnprocessed(context.Context) ([]models.Order, error)
		Update(context.Context, models.Order) error
	}

	WithdrawalRepository interface {
		Withdraw(context.Context, models.Withdrawal) error
		GetAll(context.Context) ([]models.Withdrawal, error)
	}

	AccrualWebAPI interface {
		GetOrder(context.Context, models.Order) (models.Order, error)
	}
)
