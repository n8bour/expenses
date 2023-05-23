package data

type User struct {
	ID       string `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Expenses []Expense
}
