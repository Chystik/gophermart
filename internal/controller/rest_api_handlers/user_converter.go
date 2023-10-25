package restapihandlers

import "github.com/Chystik/gophermart/internal/models"

type balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func fromDomainBalance(u models.User) balance {
	return balance{
		Current:   u.Balance,
		Withdrawn: u.Withdrawn,
	}
}
