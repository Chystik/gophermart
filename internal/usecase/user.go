package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Chystik/gophermart/internal/models"

	"github.com/avito-tech/go-transaction-manager/trm"
)

var (
	ErrNotEnoughMoney = errors.New("not enough money")
)

type userInteractor struct {
	userRepo       UserRepository
	withdrawalRepo WithdrawalRepository
	trm            trm.Manager
}

func NewUserInteractor(ur UserRepository, wr WithdrawalRepository, trm trm.Manager) *userInteractor {
	return &userInteractor{
		userRepo:       ur,
		withdrawalRepo: wr,
		trm:            trm,
	}
}

func (ui *userInteractor) Register(ctx context.Context, user models.User) error {
	err := ui.trm.Do(ctx, func(ctx context.Context) error {
		err := ui.userRepo.Create(ctx, user)
		return err
	})

	return err
}

func (ui *userInteractor) Authenticate(ctx context.Context, user models.User) error {
	actual, err := ui.userRepo.Get(ctx, user)
	if err != nil {
		return err
	}

	return user.Authenticate(actual)
}

func (ui *userInteractor) Get(ctx context.Context, user models.User) (models.User, error) {
	var err error
	err = ui.trm.Do(ctx, func(ctx context.Context) error {
		user, err = ui.userRepo.Get(ctx, user)
		return err
	})

	return user, err
}

func (ui *userInteractor) Withdraw(ctx context.Context, w models.Withdrawal, user models.User) error {
	var err error
	var actual models.User

	err = ui.trm.Do(ctx, func(ctx context.Context) error {
		actual, err = ui.userRepo.Get(ctx, user)
		if err != nil {
			return err
		}

		if actual.Balance < w.Sum {
			return ErrNotEnoughMoney
		}

		actual.Balance -= w.Sum
		actual.Withdrawn += w.Sum
		w.ProcessedAt = models.RFC3339Time{Time: time.Now()}

		err = ui.withdrawalRepo.Create(ctx, w)
		if err != nil {
			return err
		}

		return ui.userRepo.Update(ctx, actual)
	})

	return err
}

func (ui *userInteractor) GetWithdrawals(ctx context.Context) ([]models.Withdrawal, error) {
	var err error
	var w []models.Withdrawal
	err = ui.trm.Do(ctx, func(ctx context.Context) error {
		w, err = ui.withdrawalRepo.GetAll(ctx)
		return err
	})

	return w, err
}
