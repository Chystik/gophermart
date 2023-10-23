package usecase

import (
	"context"

	"github.com/Chystik/gophermart/internal/models"
)

type userInteractor struct {
	userRepo UserRepository
}

func NewUserInteractor(ur UserRepository) *userInteractor {
	return &userInteractor{
		userRepo: ur,
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
