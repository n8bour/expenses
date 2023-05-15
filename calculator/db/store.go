package db

type Storer[T any] interface {
	Insert(T) (*T, error)
	Get(int64) (*T, error)
	List() (*[]T, error)
}
