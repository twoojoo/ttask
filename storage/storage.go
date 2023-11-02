package storage

type Storage[T any] interface {
	// Push(key string, x *T, meta map[string]int64) int
	// Flush(key string) []T
	// GetSize(key string) int
	// GetAllKeys() []string
	// GetAllSizes() map[string]int
	// GetMetadata(key string) map[string]int64
	// SetMetadata(key string, meta map[string]int64) map[string]int64
	StartNewWindow(key string, elem T, start ...int64) WindowMeta
	StartNewEmptyWindow(key string, start ...int64) WindowMeta
	PushItemToWindow(k string, id string, item T) int
	GetWindowsMetadata(k string) []WindowMeta
	GetWindowMetadata(k string, id string) WindowMeta
	// SetWindowMetadata(k string, id string)
	CloseWindow(k string, id string)
	FlushWindow(k string, id string) []T
	DestroyWindow(k string, id string)
	GetWindowSize(k string, id string) int
	GetKeys() []string
}
