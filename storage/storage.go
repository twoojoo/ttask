package storage

type Storage[T any] interface {
	Push(key string, x *T, meta map[string]int64) int
	Flush(key string) []T
	GetSize(key string) int
	GetAllKeys() []string
	GetAllSizes() map[string]int
	GetMetadata(key string) map[string]int64
	SetMetadata(key string, meta map[string]int64) map[string]int64
}
