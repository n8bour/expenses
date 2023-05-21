package db

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/n8bour/expenses/calculator/data"
	"log"
	"time"
)

type UserStore struct {
	*sqlx.DB
}

func NewSqlUserStore(db *sqlx.DB) *UserStore {
	autoMigrateUser(db)
	return &UserStore{DB: db}
}

func (s *UserStore) Insert(exp data.User) (*data.User, error) {
	query := `insert into "user" (username, password) values (:username, :password) returning id`

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

func (s *UserStore) Get(id int64) (*data.User, error) {
	var r data.User
	query := `select * from "user" where id = $1`
	err := s.DB.Get(&r, query, id)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (s *UserStore) List() (*[]data.User, error) {
	var r []data.User
	query := `select * from "user"`
	err := s.DB.Select(&r, query)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func autoMigrateUser(db *sqlx.DB) {
	query := `create table if not exists "user"(id serial primary key,username varchar not null,password varchar not null);`

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	_, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Print("ERROR CREATING TABLE: ", err)
	}
}
