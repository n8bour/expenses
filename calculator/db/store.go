package db

import "context"

type Storer[T any] interface {
	Insert(context.Context, T) (*T, error)
	Get(context.Context, string) (*T, error)
	List(context.Context) (*[]T, error)
}
