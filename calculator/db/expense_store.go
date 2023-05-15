package db

import (
	"context"
	"database/sql"
	"github.com/n8bour/expenses/calculator/data"
	"log"
	"time"
)

type ExpenseStore struct {
	*sql.DB
}

func NewSqlExpenseStore(db *sql.DB) *ExpenseStore {
	setup(db)
	return &ExpenseStore{DB: db}
}

func (s *ExpenseStore) Insert(exp data.Expense) (*data.Expense, error) {
	query := "insert into expense (type, value) values ($1,$2) returning id"

	r := s.QueryRow(query, exp.Type, exp.Value)

	var id int64
	err := r.Scan(&id)
	if err != nil {
		return nil, err
	}
	exp.ID = id

	return &exp, nil
}

func (s *ExpenseStore) Get(id int64) (*data.Expense, error) {
	var r data.Expense
	query := "select * from expense where id = $1"
	row := s.QueryRow(query, id)
	err := row.Scan(&r.ID, &r.Type, &r.Value)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (s *ExpenseStore) List() (*[]data.Expense, error) {
	//TODO implement me
	panic("implement me")
}

func setup(db *sql.DB) {
	query := `create table if not exists expense(id serial primary key,type varchar not null,value float4 not null);`

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	_, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Print("ERROR CREATING TABLE: ", err)
	}
}
