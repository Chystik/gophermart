package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Chystik/gophermart/internal/models"
)

var (
	ErrNotEnoughMoney = errors.New("not enough money")
)

type userInteractor struct {
	userRepo       UserRepository
	withdrawalRepo WithdrawalRepository
}

func NewUserInteractor(ur UserRepository, wr WithdrawalRepository) *userInteractor {
	return &userInteractor{
		userRepo:       ur,
		withdrawalRepo: wr,
	}
}

func (ui *userInteractor) Register(ctx context.Context, user models.User) error {
	return ui.userRepo.Create(ctx, user)
}

func (ui *userInteractor) Authenticate(ctx context.Context, user models.User) error {
	actual, err := ui.userRepo.Get(ctx, user)
	if err != nil {
		return err
	}

	return user.Authenticate(actual)
}

func (ui *userInteractor) Get(ctx context.Context, user models.User) (models.User, error) {
	return ui.userRepo.Get(ctx, user)
}

func (ui *userInteractor) Withdraw(ctx context.Context, w models.Withdrawal, user models.User) error {
	actual, err := ui.userRepo.Get(ctx, user)
	if err != nil {
		return err
	}

	if actual.Balance < w.Sum {
		return ErrNotEnoughMoney
	}

	actual.Balance -= w.Sum
	actual.Withdrawn += w.Sum
	w.ProcessedAt = models.RFC3339Time{Time: time.Now()}

	err = ui.withdrawalRepo.Withdraw(ctx, w)
	if err != nil {
		return err
	}

	return ui.userRepo.Update(ctx, actual)
}

func (ui *userInteractor) GetWithdrawals(ctx context.Context) ([]models.Withdrawal, error) {
	return ui.withdrawalRepo.GetAll(ctx)
}
