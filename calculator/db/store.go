package db

import (
	"context"

	"github.com/n8bour/expenses/calculator/data"
)

type Storer[T any] interface {
	Insert(context.Context, T) (*T, error)
	Get(context.Context, string) (*T, error)
	List(context.Context) (*[]T, error)
}

type UserStorer interface {
	Storer[data.User]
	GetByUsername(ctx context.Context, username string) (*data.User, error)
	GetCurrentUser(ctx context.Context) (*data.User, error)
}
