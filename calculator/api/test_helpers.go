package api

import (
	"context"
	"os/exec"

	"github.com/n8bour/expenses/calculator/data"
	"golang.org/x/crypto/bcrypt"
)

type UserMemoryDB struct {
	store map[string]data.User
}

func (s UserMemoryDB) GetByUsername(_ context.Context, username string) (*data.User, error) {
	for _, v := range s.store {
		if v.Username == username {
			return &v, nil
		}
	}
	return nil, NotResourceNotFound(username)
}

func (s UserMemoryDB) GetCurrentUser(ctx context.Context) (*data.User, error) {
	username := ctx.Value("currentUser").(map[string]string)["username"]
	for _, v := range s.store {
		if v.Username == username {
			return &v, nil
		}
	}
	return nil, NotResourceNotFound(username)
}

func NewUserMemoryDB() *UserMemoryDB {
	return &UserMemoryDB{store: map[string]data.User{}}
}

func (s UserMemoryDB) Insert(_ context.Context, u data.User) (*data.User, error) {
	uuid, _ := exec.Command("uuidgen").Output()
	u.ID = string(uuid)
	pw, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 4)
	u.Password = string(pw)
	s.store[u.ID] = u

	return &u, nil
}

func (s UserMemoryDB) Get(_ context.Context, id string) (*data.User, error) {
	user, ok := s.store[id]
	if !ok {
		return nil, NotResourceNotFound(id)
	}

	return &user, nil
}

func (s UserMemoryDB) List(_ context.Context) (*[]data.User, error) {
	var users []data.User
	for _, v := range s.store {
		users = append(users, v)
	}

	return &users, nil
}

type ExpenseMemoryDB struct {
	store map[string]data.Expense
}

func NewExpenseMemoryDB() *ExpenseMemoryDB {
	return &ExpenseMemoryDB{store: map[string]data.Expense{}}
}

func (e ExpenseMemoryDB) Insert(_ context.Context, exp data.Expense) (*data.Expense, error) {
	uuid, _ := exec.Command("uuidgen").Output()
	exp.ID = string(uuid)
	e.store[exp.ID] = exp

	return &exp, nil
}

func (e ExpenseMemoryDB) Get(_ context.Context, id string) (*data.Expense, error) {
	exp, ok := e.store[id]
	if !ok {
		return nil, NotResourceNotFound(id)
	}

	return &exp, nil
}

func (e ExpenseMemoryDB) List(_ context.Context) (*[]data.Expense, error) {
	var expenses []data.Expense
	for _, v := range e.store {
		expenses = append(expenses, v)
	}

	return &expenses, nil
}
