package db

type Storer[T any] interface {
	Insert(T) (*T, error)
	Get(string2 string) (*T, error)
	List() (*[]T, error)
}
