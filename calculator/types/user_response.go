package types

type UserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Expenses []ExpenseRequest
}
