package db

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/n8bour/expenses/calculator/data"
	"log"
	"time"
)

type ExpenseStore struct {
	*sqlx.DB
}

func NewSqlExpenseStore(db *sqlx.DB) *ExpenseStore {
	autoMigrateExpense(db)
	return &ExpenseStore{DB: db}
}

func (s *ExpenseStore) Insert(exp data.Expense) (*data.Expense, error) {
	query := "insert into expense (type, value) values (:type, :value) returning id"

	//r := s.QueryRow(query, exp.Type, exp.Value)

	row, err := s.NamedQuery(query, &exp)
	if err != nil {
		return nil, err
	}
	row.Next()
	err = row.StructScan(&exp)
	if err != nil {
		return nil, err
	}

	return &exp, nil
}

func (s *ExpenseStore) Get(id int64) (*data.Expense, error) {
	var r data.Expense
	query := "select * from expense where id = $1"
	err := s.DB.Get(&r, query, id)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (s *ExpenseStore) List() (*[]data.Expense, error) {
	var r []data.Expense
	query := "select * from expense"
	err := s.DB.Select(&r, query)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func autoMigrateExpense(db *sqlx.DB) {
	query := `create table if not exists expense(id serial primary key,type varchar not null,value float4 not null);`

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	_, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Print("ERROR CREATING TABLE: ", err)
	}
}
