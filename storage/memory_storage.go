package storage

import (
	"sync"

	"github.com/twoojoo/ttask/task"
)

func Memory[T any]() *MemoryStorage[task.Message[T]] {
	return &MemoryStorage[task.Message[T]]{
		metadata: map[string]map[string]int64{},
		mem:      map[string][]task.Message[T]{},
	}
}

type MemoryStorage[T any] struct {
	metadata map[string]map[string]int64
	mem      map[string][]T
	mu       sync.Mutex
}

func (s *MemoryStorage[T]) GetAllKeys() []string {
	keys := []string{}

	for k := range s.mem {
		keys = append(keys, k)
	}

	return keys
}

func (s *MemoryStorage[T]) GetAllSizes() map[string]int {
	sizes := map[string]int{}

	for k, v := range s.mem {
		sizes[k] = len(v)
	}

	return sizes
}

func (s *MemoryStorage[T]) GetSize(key string) int {
	return len(s.mem[key])
}

func (s *MemoryStorage[T]) GetMetadata(key string) map[string]int64 {
	if s.metadata[key] == nil {
		s.metadata[key] = map[string]int64{}
	}

	return s.metadata[key]
}

func (s *MemoryStorage[T]) SetMetadata(key string, meta map[string]int64) map[string]int64 {
	if s.metadata[key] == nil {
		s.metadata[key] = map[string]int64{}
	}

	for k, v := range meta {
		s.metadata[key][k] = v
	}

	return s.metadata[key]
}

func (s *MemoryStorage[T]) Push(key string, item *T, meta map[string]int64) int {
	s.mu.Lock()

	if s.metadata[key] == nil {
		s.metadata[key] = map[string]int64{}
	}

	for k, v := range meta {
		s.metadata[key][k] = v
	}

	s.mem[key] = append(s.mem[key], *item)

	s.mu.Unlock()

	return len(s.mem[key])
}

func (s *MemoryStorage[T]) Flush(key string) []T {
	result := []T{}

	s.mu.Lock()

	for _, v := range s.mem[key] {
		result = append(result, v)
	}

	s.mem[key] = []T{}

	s.mu.Unlock()

	return result
}
