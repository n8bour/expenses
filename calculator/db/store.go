package db

type Storer[T any] interface {
	Insert(T) (*T, error)
	Get(id string) (*T, error)
	List() (*[]T, error)
}
