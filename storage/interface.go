package storage

type StorageInterface[T any] struct {
	storage *Storage[T]
}

func NewStorageInterface[T any](s *Storage[T]) StorageInterface[T]{
	return StorageInterface[T]{storage: s}
}

func (s *StorageInterface[T]) GetKeys() []string {
	return (*s.storage).GetKeys()
}

func (s *StorageInterface[T]) StartNewEmptyWindow(key string, start ...int64) WindowMeta {
	return (*s.storage).StartNewEmptyWindow(key, start...)
}

func (s *StorageInterface[T]) StartNewWindow(key string, elem T, start ...int64) WindowMeta {
	return (*s.storage).StartNewWindow(key, elem, start...)
}

func (s *StorageInterface[T]) GetWindowSize(k string, id string) int {
	return (*s.storage).GetWindowSize(k, id)
}

func (s *StorageInterface[T]) GetWindowsMetadata(k string) []WindowMeta {
	return (*s.storage).GetWindowsMetadata(k)
}

func (s *StorageInterface[T]) GetWindowMetadata(k string, id string) WindowMeta {
	return (*s.storage).GetWindowMetadata(k, id)
}

func (s *StorageInterface[T]) CloseWindow(k string, id string) {
	(*s.storage).CloseWindow(k, id)
}

func (s *StorageInterface[T]) FlushWindow(k string, id string) []T {
	return (*s.storage).FlushWindow(k, id)
}

func (s *StorageInterface[T]) PushItemToWindow(k string, id string, item T) int {
	return (*s.storage).PushItemToWindow(k, id, item)
}

func (s *StorageInterface[T]) DestroyWindow(k string, id string) {
	(*s.storage).DestroyWindow(k, id)
}
