package db

import (
	"context"
	"database/sql"
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
	query := "insert into expense (type, value, user_id) values (:type, :value, :user_id) returning id"

	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()

	tx, err := s.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	row, err := tx.NamedQuery(query, &exp)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	row.Next()
	err = row.StructScan(&exp)
	if err != nil {
		return nil, err
	}

	return &exp, tx.Commit()
}

func (s *ExpenseStore) Get(id string) (*data.Expense, error) {
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
	query := `create table if not exists expense(id uuid default gen_random_uuid() primary key,type varchar not null,value float4 not null, user_id uuid not null, constraint fk_user foreign key(user_id) references "user"(id));`

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	_, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Fatal("ERROR CREATING TABLE: ", err)
	}

	view := `create or replace view user_expenses as
    select id, username, password,
           (select array_to_json(array_agg(row_to_json(expenses_list.*))) as array_to_json
            from (select e.id, e.type, e.value from expense e where user_id = "user".id) expenses_list) as expenses
from "user";`

	_, err = db.ExecContext(ctx, view)
	if err != nil {
		log.Fatal("ERROR CREATING VIEW FOR ONE TO MANY RELATION: ", err)
	}
}
