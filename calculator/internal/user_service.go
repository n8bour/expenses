package internal

import (
	"github.com/n8bour/expenses/calculator/data"
	"github.com/n8bour/expenses/calculator/db"
	"github.com/n8bour/expenses/calculator/types"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	store db.Storer[data.User]
}

func NewUserService(store db.Storer[data.User]) *UserService {
	return &UserService{
		store: store,
	}
}

func (s *UserService) CreateUser(usr types.UserRequest) (*types.UserResponse, error) {

	password, err := bcrypt.GenerateFromPassword([]byte(usr.Password), 4)
	if err != nil {
		return nil, err
	}

	//usr.Password = password

	err = bcrypt.CompareHashAndPassword(password, []byte(usr.Password))
	if err != nil {
		return nil, err
	}

	/*r, err := s.store.Insert(usr.ToUser())
	if err != nil {
		return nil, err
	}

	return &types.UserResponse{
		ID:       r.ID,
		Username: r.Username,
	}, nil*7
	*/

	return nil, nil
}

func (s *UserService) GetUser(id string) (result *types.UserResponse, err error) {

	//return result.FromUser(r), nil
	return nil, nil
}

func (s *UserService) ListUsers() (*[]types.UserResponse, error) {
	return nil, nil
}
