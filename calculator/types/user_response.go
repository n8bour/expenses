package types

import "github.com/n8bour/expenses/calculator/data"

type UserResponse struct {
	ID       string            `json:"id"`
	Username string            `json:"username"`
	Email    string            `json:"email"`
	Expenses []ExpenseResponse `json:"expenses,omitempty"`
}

func (ur UserResponse) FromUser(u data.User) *UserResponse {
	var expenses []ExpenseResponse
	if len(u.Expenses) > 0 {
		for _, es := range u.Expenses {
			exr := ExpenseResponse{}
			expenses = append(expenses, *exr.FromExpense(es))
		}
	}

	return &UserResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Expenses: expenses,
	}
}
