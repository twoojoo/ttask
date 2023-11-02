package storage

type Storage[T any] interface {
	// Push(key string, x *T, meta map[string]int64) int
	// Flush(key string) []T
	// GetSize(key string) int
	// GetAllKeys() []string
	// GetAllSizes() map[string]int
	// GetMetadata(key string) map[string]int64
	// SetMetadata(key string, meta map[string]int64) map[string]int64
	StartNewWindow(key string, md map[string]int64, elems ...T) WindowMeta
	PushItemToWindow(k string, id string, item T, md map[string]int64) int
	GetWindowsMetadata(k string) []WindowMeta
	SetWindowMetadata(k string, id string, md map[string]int64)
	CloseWindow(k string, id string) []T
	DestroyWindow(k string, id string)
	GetWindowSize(k string, id string) int
}
