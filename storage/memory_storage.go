package storage

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/twoojoo/ttask/task"
)

func Memory[T any]() *MemoryStorage[task.Message[T]] {
	return &MemoryStorage[task.Message[T]]{
		windows: map[string][]window[task.Message[T]]{},
	}
}

type MemoryStorage[T any] struct {
	windows map[string][]window[T]
	mu      sync.Mutex
}

func (s *MemoryStorage[T]) StartNewWindow(key string, md map[string]int64, elems ...T) WindowMeta {
	if s.windows[key] == nil {
		s.windows[key] = []window[T]{}
	}

	newMd := map[string]int64{
		"start": time.Now().UnixMicro(),
	}

	winMeta := WindowMeta{
		Id:       uuid.NewString(),
		Metadata: mergeMetadata(newMd, md),
	}

	s.windows[key] = append(s.windows[key], window[T]{
		id:       winMeta.Id,
		metadata: winMeta.Metadata,
		elems:    elems,
	})

	return winMeta
}

func (s *MemoryStorage[T]) GetWindowSize(k string, id string) int {
	if s.windows[k] == nil {
		return 0
	}

	for _, win := range s.windows[k] {
		if win.id == id {
			return len(win.elems)
		}
	}

	return 0
}

func (s *MemoryStorage[T]) GetWindowsMetadata(k string) []WindowMeta {
	if s.windows[k] == nil {
		return []WindowMeta{}
	}

	meta := []WindowMeta{}

	for i := 0; i < len(s.windows[k]); i++ {
		wm := WindowMeta{
			Id:       s.windows[k][i].id,
			Metadata: s.windows[k][i].metadata,
		}

		meta = append(meta, wm)
	}

	return meta
}

func (s *MemoryStorage[T]) SetWindowMetadata(k string, id string, md map[string]int64) {
	if s.windows[k] == nil {
		return
	}

	for _, win := range s.windows[k] {
		if win.id == id {
			win.metadata = mergeMetadata(win.metadata, md)
		}
	}
}

func (s *MemoryStorage[T]) CloseWindow(k string, id string) []T {
	if s.windows[k] == nil {
		return []T{}
	}

	elems := []T{}

	for i := range s.windows[k] {
		if s.windows[k][i].id == id {
			elems = s.windows[k][i].elems
			s.windows[k] = append(s.windows[k][:i], s.windows[k][i+1:]...)
			return elems
		}
	}

	return elems
}

func (s *MemoryStorage[T]) PushItemToWindow(k string, id string, item T, md map[string]int64) int {
	if s.windows[k] == nil {
		s.windows[k] = []window[T]{}
	}

	for i, _ := range s.windows[k] {
		if s.windows[k][i].id == id {
			s.windows[k][i].elems = append(s.windows[k][i].elems, item)
			s.windows[k][i].metadata = mergeMetadata(s.windows[k][i].metadata, md)
		}

		return len(s.windows[k][i].elems)
	}

	panic("a window should exist for this id and key")
}

func (s *MemoryStorage[T]) DestroyWindow(k string, id string) {
	if s.windows[k] == nil {
		return
	}

	for i, win := range s.windows[k] {
		if win.id == id {
			s.windows[k] = append(s.windows[k][:i], s.windows[k][i+1:]...)
			break
		}
	}
}

// func (s *MemoryStorage[T]) GetAllKeys() []string {
// 	keys := []string{}

// 	for k := range s.windows {
// 		winKeys := []string{}
// 		for k1 := range s.windows[k].elems {
// 			winKeys = append(winKeys, k1)
// 		}

// 		keys = append(keys, k)
// 	}

// 	return keys
// }

// func (s *MemoryStorage[T]) GetAllSizes() map[string]int {
// 	sizes := map[string]int{}

// 	for k, v := range s.mem {
// 		sizes[k] = len(v)
// 	}

// 	return sizes
// }

// func (s *MemoryStorage[T]) GetSize(key string) int {
// 	return len(s.mem[key])
// }

// func (s *MemoryStorage[T]) GetMetadata(key string) map[string]int64 {
// 	if s.metadata[key] == nil {
// 		s.metadata[key] = map[string]int64{}
// 	}

// 	return s.metadata[key]
// }

// func (s *MemoryStorage[T]) SetMetadata(key string, meta map[string]int64) map[string]int64 {
// 	if s.metadata[key] == nil {
// 		s.metadata[key] = map[string]int64{}
// 	}

// 	for k, v := range meta {
// 		s.metadata[key][k] = v
// 	}

// 	return s.metadata[key]
// }

// func (s *MemoryStorage[T]) Push(key string, item *T, meta map[string]int64) int {
// 	s.mu.Lock()

// 	if s.metadata[key] == nil {
// 		s.metadata[key] = map[string]int64{}
// 	}

// 	for k, v := range meta {
// 		s.metadata[key][k] = v
// 	}

// 	s.mem[key] = append(s.mem[key], *item)

// 	s.mu.Unlock()

// 	return len(s.mem[key])
// }

// func (s *MemoryStorage[T]) Flush(key string) []T {
// 	result := []T{}

// 	s.mu.Lock()

// 	for _, v := range s.mem[key] {
// 		result = append(result, v)
// 	}

// 	s.mem[key] = []T{}

// 	s.mu.Unlock()

// 	return result
// }
