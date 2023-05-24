package data

import "time"

type Expense struct {
	ID        string    `db:"id"`
	Type      string    `db:"type"`
	Value     float32   `db:"value"`
	CreatedAt time.Time `db:"created_at"`
	UserID    string    `db:"user_id"`
}
