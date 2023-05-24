package db

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/n8bour/expenses/calculator/data"
)

type UserStore struct {
	*sqlx.DB
}

func NewSQLUserStore(db *sqlx.DB) *UserStore {
	autoMigrateUser(db)
	return &UserStore{DB: db}
}

func (s *UserStore) Insert(ctx context.Context, exp data.User) (*data.User, error) {
	query := `insert into "user" (username, email ,password) values (:username, :email,:password) returning id`

	row, err := s.NamedQueryContext(ctx, query, &exp)
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

func (s *UserStore) Get(ctx context.Context, id string) (*data.User, error) {
	var (
		r     data.User
		rjson json.RawMessage
	)
	query := `select row_to_json(row) from (select * from user_expenses) row where row.id = $1`

	err := s.DB.GetContext(ctx, &rjson, query, id)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(rjson, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (s *UserStore) GetCurrentUser(ctx context.Context) (*data.User, error) {
	var r data.User

	query := `select * from "user" u where username = $1`

	username := ctx.Value("currentUser").(map[string]string)["username"]

	err := s.DB.GetContext(ctx, &r, query, username)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (s *UserStore) GetByUsername(ctx context.Context, username string) (*data.User, error) {
	var r data.User
	query := `select * from "user" u where username = $1`

	err := s.DB.GetContext(ctx, &r, query, username)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (s *UserStore) List(ctx context.Context) (*[]data.User, error) {
	var (
		r     []data.User
		rjson []json.RawMessage
	)
	query := `with users_exp as (select u.*, array_to_json(array_agg( row_to_json(e.*))) as expenses from "user" u left join expense e on u.id = e.user_id group by u.id)
select row_to_json(users_exp)
from users_exp`

	err := s.DB.SelectContext(ctx, &rjson, query)
	if err != nil {
		return nil, err
	}

	for _, message := range rjson {
		u := data.User{}
		err = json.Unmarshal(message, &u)
		if err != nil {
			return nil, err
		}
		if u.Expenses[0].Type == "" {
			u.Expenses = nil
		}
		r = append(r, u)
	}

	return &r, nil
}

func autoMigrateUser(db *sqlx.DB) {
	query := `create table if not exists "user"(id uuid default gen_random_uuid() primary key,username varchar not null, email varchar not null, password varchar not null);`

	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()

	_, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Panicf("ERROR CREATING TABLE: %s", err)
	}
}
