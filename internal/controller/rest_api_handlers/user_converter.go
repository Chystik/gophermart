package restapihandlers

import "github.com/Chystik/gophermart/internal/models"

type balance struct {
	Current   models.Money `json:"current"`
	Withdrawn models.Money `json:"withdrawn"`
}

func fromDomainBalance(u models.User) balance {
	return balance{
		Current:   u.Balance,
		Withdrawn: u.Withdrawn,
	}
}
