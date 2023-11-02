package storage

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/twoojoo/ttask/task"
)

func Redis[T any](id string, redis *redis.Client) *RedisStorage[task.Message[T]] {
	return &RedisStorage[task.Message[T]]{
		client: redis,
		id: id,
	}
}

type RedisStorage[T any] struct {
	id     string
	client *redis.Client
}

func (s RedisStorage[T]) redisKeysK(k string) string {
	return s.id + ".k." + k
}

func (s RedisStorage[T]) redisWinK(k string, id string) string {
	return s.id + ".w." + k + "." + id
}

func (s *RedisStorage[T]) GetKeys() []string {
	keys := []string{}

	s.client.Keys(context.Background(), s.id+".k.*")

	return keys
}

func (s *RedisStorage[T]) StartNewEmptyWindow(key string, start ...int64) WindowMeta {
	// if s.windows[key] == nil {
	// 	s.windows[key] = make(map[string]window[T], 0)
	// }

	// id := uuid.NewString()

	winMeta := WindowMeta{}
	// 	Id:       id,
	// 	Metadata: map[string]int64{},
	// 	Start:    time.Now().UnixMilli(),
	// }

	// if len(start) > 0 {
	// 	winMeta.Start = start[0]
	// }

	// s.windows[key][id] = window[T]{
	// 	id:       winMeta.Id,
	// 	metadata: winMeta.Metadata,
	// 	elems:    []T{},
	// 	start:    winMeta.Start,
	// 	end:      winMeta.End,
	// }

	return winMeta
}

func (s *RedisStorage[T]) StartNewWindow(key string, elem T, start ...int64) WindowMeta {
	// if s.windows[key] == nil {
	// 	s.windows[key] = map[string]window[T]{}
	// }

	// id := uuid.NewString()

	winMeta := WindowMeta{}
	// }

	// if len(start) > 0 {
	// 	winMeta.Start = start[0]
	// }

	// s.windows[key][id] = window[T]{
	// 	id:       winMeta.Id,
	// 	metadata: winMeta.Metadata,
	// 	elems:    []T{elem},
	// 	start:    winMeta.Start,
	// 	end:      winMeta.End,
	// 	last:     winMeta.Last,
	// }

	return winMeta
}

func (s *RedisStorage[T]) GetWindowSize(k string, id string) int {
	// if s.windows[k] == nil {
	// 	return 0
	// }

	// for _, win := range s.windows[k] {
	// 	if win.id == id {
	// 		return len(win.elems)
	// 	}
	// }

	return 0
}

func (s *RedisStorage[T]) GetWindowsMetadata(k string) []WindowMeta {
	// if s.windows[k] == nil {
	// 	return []WindowMeta{}
	// }

	meta := []WindowMeta{}

	// for id := range s.windows[k] {
	// 	wm := WindowMeta{
	// 		Id:       s.windows[k][id].id,
	// 		Metadata: s.windows[k][id].metadata,
	// 		Start:    s.windows[k][id].start,
	// 		End:      s.windows[k][id].end,
	// 		Last:     s.windows[k][id].last,
	// 	}

	// 	meta = append(meta, wm)
	// }

	return meta
}

func (s *RedisStorage[T]) GetWindowMetadata(k string, id string) WindowMeta {
	// if winById, ok := s.windows[k]; ok {
	// 	if win, ok := winById[id]; ok {
	// 		return WindowMeta{
	// 			Id:       win.id,
	// 			Last:     win.last,
	// 			End:      win.end,
	// 			Start:    win.start,
	// 			Metadata: win.metadata,
	// 		}
	// 	}
	// }

	return WindowMeta{}
}

func (s *RedisStorage[T]) CloseKeyWindows(k string, id string) {
	// if winById, ok := s.windows[k]; ok {
	// 	if win, ok := winById[id]; ok {
	// 		win.end = time.Now().UnixMilli()
	// 	}
	// }
}

func (s *RedisStorage[T]) CloseWindow(k string, id string) {
	// if winById, ok := s.windows[k]; ok {
	// 	if win, ok := winById[id]; ok {
	// 		win.end = time.Now().UnixMilli()
	// 		s.windows[k][id] = win
	// 	}
	// }
}

func (s *RedisStorage[T]) FlushWindow(k string, id string) []T {
	// if s.windows[k] == nil {
	// 	return []T{}
	// }

	// if winById, ok := s.windows[k]; ok {
	// 	if win, ok := winById[id]; ok {
	// 		delete(s.windows[k], id)
	// 		return win.elems
	// 	}
	// }

	return []T{}
}

func (s *RedisStorage[T]) PushItemToWindow(k string, id string, item T) int {
	size := 0

	// if winById, ok := s.windows[k]; ok {
	// 	if win, ok := winById[id]; ok {
	// 		win.elems = append(win.elems, item)
	// 		s.windows[k][id] = win
	// 		size = len(win.elems)
	// 	}
	// }

	return size
}

func (s *RedisStorage[T]) DestroyWindow(k string, id string) {
	// if winById, ok := s.windows[k]; ok {
	// 	delete(winById, id)
	// }
}
