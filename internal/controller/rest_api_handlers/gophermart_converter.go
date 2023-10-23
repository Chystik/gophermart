package restapihandlers

import "github.com/Chystik/gophermart/internal/models"

type credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type order struct {
	Number int
	User   string
}

func toDomainUser(c credentials) models.User {
	return models.User{
		Login:    c.Login,
		Password: c.Password,
	}
}

func toDomainOrder(o order) models.Order {
	return models.Order{
		Number: o.Number,
		User:   o.User,
	}
}
