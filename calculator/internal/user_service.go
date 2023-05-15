package internal

import (
	"github.com/n8bour/expenses/calculator/data"
	"github.com/n8bour/expenses/calculator/db"
)

type UserService struct {
	store db.Storer[data.User]
}

func NewUserService(store db.Storer[data.User]) *UserService {
	return &UserService{
		store: store,
	}
}
