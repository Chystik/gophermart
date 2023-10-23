package usecase

import (
	"context"
	"errors"

	"github.com/Chystik/gophermart/internal/models"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrWrongCreds = errors.New("mismatch login or password")
)

type gophermartInteractor struct {
	userRepo  UserRepository
	orderRepo OrderRepository
}

func NewGophermartInteractor(ur UserRepository, or OrderRepository) *gophermartInteractor {
	return &gophermartInteractor{
		userRepo:  ur,
		orderRepo: or,
	}
}

func (gm *gophermartInteractor) RegisterUser(ctx context.Context, user models.User) error {
	return gm.userRepo.Create(ctx, user)
}

func (gm *gophermartInteractor) AuthenticateUser(ctx context.Context, user models.User) error {
	actual, err := gm.userRepo.Get(ctx, user)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(actual.Password), []byte(user.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return ErrWrongCreds
		}
		return err
	}

	return nil
}

func (gm *gophermartInteractor) GetUser(ctx context.Context, user models.User) (models.User, error) {
	return gm.userRepo.Get(ctx, user)
}

func (gm *gophermartInteractor) CreateOrder(ctx context.Context, order models.Order) error {
	return gm.orderRepo.Create(ctx, order)
}

func (gm *gophermartInteractor) GetOrders(ctx context.Context) ([]models.Order, error) {
	return gm.orderRepo.GetList(ctx)
}
