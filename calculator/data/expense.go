package data

type Expense struct {
	ID    int64   `db:"id"`
	Type  string  `db:"type"`
	Value float32 `db:"value"`
	//UserID int64
}
