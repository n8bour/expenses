package types

import "github.com/n8bour/expenses/calculator/data"

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (ur UserRequest) ToUser() data.User {
	return data.User{
		Username: ur.Username,
		Password: ur.Password,
	}

}