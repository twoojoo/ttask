package storage

import "github.com/twoojoo/ttask/task"

func Memory[T any]() *MemoryStorage[task.Message[T]] {
	return &MemoryStorage[task.Message[T]]{
		mem: map[string][]task.Message[T]{},
	}
}

type MemoryStorage[T any] struct {
	mem map[string][]T
}

// func (s *MemoryStorage[T]) Init() {
// 	s.mem = map[string][]T{}
// }

func (s *MemoryStorage[T]) Push(key string, item *T) int {
	s.mem[key] = append(s.mem[key], *item)
	return len(s.mem[key])
}

func (s *MemoryStorage[T]) Flush(key string) []T {
	result := []T{}

	for _, v := range s.mem[key] {
		result = append(result, v)
	}

	s.mem[key] = []T{}

	return result
}
