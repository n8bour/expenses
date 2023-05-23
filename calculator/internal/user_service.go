package internal

import (
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

func (s *UserService) CreateUser(usr types.UserRequest) (*types.UserResponse, error) {

	password, err := bcrypt.GenerateFromPassword([]byte(usr.Password), 4)
	if err != nil {
		return nil, err
	}

	usr.Password = string(password)

	r, err := s.store.Insert(usr.ToUser())
	if err != nil {
		return nil, err
	}

	return &types.UserResponse{
		ID:       r.ID,
		Username: r.Username,
	}, nil
}

func (s *UserService) GetUser(id string) (*types.UserResponse, error) {
	var result types.UserResponse

	u, err := s.store.Get(id)
	if err != nil {
		return nil, err
	}
	return result.FromUser(u), nil
}

func (s *UserService) ListUsers() (*[]types.UserResponse, error) {
	list, err := s.store.List()
	if err != nil {
		return nil, err
	}

	var result []types.UserResponse
	for _, user := range *list {
		u := types.UserResponse{}
		result = append(result, *u.FromUser(&user))
	}

	return &result, nil
}
