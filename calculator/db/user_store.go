package db

import (
	"context"
	"database/sql"
	"github.com/n8bour/expenses/calculator/data"
	"log"
	"time"
)

type UserStore struct {
	*sql.DB
}

func NewSqlUserStore(db *sql.DB) *UserStore {
	return &UserStore{DB: db}
}
func (s *UserStore) Post(exp data.User) (*data.User, error) {
	/*query := "INSERT INTO user (username, password) VALUES (?, ?)"
	r, err := s.Exec(query, exp.Type, exp.Value)
	if err != nil {
		return nil, err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return nil, err
	}
	exp.ID = id*/

	return &exp, nil
}

func (s *UserStore) init() {
	query := `CREATE TABLE IF NOT EXISTS user(id int primary key auto_increment, "username" varchar,  "password" varchar)`
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()

	res, err := s.ExecContext(ctx, query)
	if err != nil {
		log.Print("ERROR CREATING TABLE: ", err)
	}

	log.Println("Table Successfully Created", res)

}
