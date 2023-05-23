package data

type Expense struct {
	ID     string  `db:"id"`
	Type   string  `db:"type"`
	Value  float32 `db:"value"`
	UserID string  `db:"user_id"`
}
