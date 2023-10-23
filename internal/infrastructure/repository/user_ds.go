package repository

import (
	"github.com/Chystik/gophermart/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type dsUser struct {
	Login     string  `db:"login"`
	Password  string  `db:"password"`
	Balance   float64 `db:"balance"`
	Withdrawn float64 `db:"withdrawn"`
}

func fromDomainUser(u models.User) (dsUser, error) {
	hash, err := fromDomainPassword(u.Password)
	if err != nil {
		return dsUser{}, err
	}

	return dsUser{
		Login:     u.Login,
		Password:  hash,
		Balance:   u.Balance,
		Withdrawn: u.Withdrawn,
	}, nil
}

func fromDomainPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 8)
	return string(hash), err
}

/* func comparePassword(dsU dsUser, u models.User) error {
	return bcrypt.CompareHashAndPassword([]byte(dsU.Password), []byte(u.Password))
} */
