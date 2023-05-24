package internal

import (
	"context"

	"github.com/n8bour/expenses/calculator/data"
	"github.com/n8bour/expenses/calculator/db"
	"github.com/n8bour/expenses/calculator/types"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	store        db.Storer[data.User]
	expenseStore db.Storer[data.Expense]
}

func NewUserService(store db.Storer[data.User], es db.Storer[data.Expense]) *UserService {
	return &UserService{
		store:        store,
		expenseStore: es,
	}
}

func (s *UserService) CreateUser(ctx context.Context, usr types.UserRequest) (*types.UserResponse, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(usr.Password), 4)
	if err != nil {
		return nil, err
	}

	usr.Password = string(password)

	r, err := s.store.Insert(ctx, usr.ToUser())
	if err != nil {
		return nil, err
	}

	resp := types.UserResponse{}.FromUser(*r)
	return resp, nil
}

func (s *UserService) GetUser(ctx context.Context, id string) (*types.UserResponse, error) {
	var result types.UserResponse

	u, err := s.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return result.FromUser(*u), nil
}

func (s *UserService) ListUsers(ctx context.Context) (*[]types.UserResponse, error) {
	list, err := s.store.List(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]types.UserResponse, 0)
	for _, user := range *list {
		u := types.UserResponse{}
		result = append(result, *u.FromUser(user))
	}

	return &result, nil
}
