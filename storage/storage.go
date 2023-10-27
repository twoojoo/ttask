package storage

type Storage[T any] interface {
	Push(key string, x *T) int
	// Pop(key string, n int) []*T
	Flush(key string) []T
	// Size(key string) int
}
