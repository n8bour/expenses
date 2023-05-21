package types

import "github.com/n8bour/expenses/calculator/data"

type UserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Expenses []ExpenseRequest
}

func (ur UserResponse) FromUser(u *data.User) *UserResponse {
	return &UserResponse{
		ID:       u.ID,
		Username: u.Username,
	}
}
