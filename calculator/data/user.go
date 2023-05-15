package data

type User struct {
	ID       int64
	Username string
	Password string
	Expenses []Expense
}
