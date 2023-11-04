package storage

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

func Memory[T any]() *MemoryStorage[Message[T]] {
	return &MemoryStorage[Message[T]]{
		windows: map[string]map[string]window[Message[T]]{},
	}
}

type Checkpoint[T any] struct {
	id      string
	message T
}

type MemoryStorage[T any] struct {
	windows     map[string]map[string]window[T]
	checkpoints map[string]Checkpoint[T]
	mu          sync.Mutex
}

func (s *MemoryStorage[T]) GetKeys() []string {
	keys := []string{}

	for k := range s.windows {
		keys = append(keys, k)
	}

	return keys
}

func (s *MemoryStorage[T]) StartNewEmptyWindow(key string, start ...int64) WindowMeta {
	if s.windows[key] == nil {
		s.windows[key] = make(map[string]window[T], 0)
	}

	id := uuid.NewString()

	winMeta := WindowMeta{
		Id:       id,
		Metadata: map[string]int64{},
		Start:    time.Now().UnixMilli(),
	}

	if len(start) > 0 {
		winMeta.Start = start[0]
	}

	s.windows[key][id] = window[T]{
		id:       winMeta.Id,
		metadata: winMeta.Metadata,
		elems:    []T{},
		start:    winMeta.Start,
		end:      winMeta.End,
	}

	return winMeta
}

func (s *MemoryStorage[T]) StartNewWindow(key string, elem T, start ...int64) WindowMeta {
	if s.windows[key] == nil {
		s.windows[key] = map[string]window[T]{}
	}

	id := uuid.NewString()

	winMeta := WindowMeta{
		Id:       id,
		Metadata: map[string]int64{},
		Start:    time.Now().UnixMilli(),
		Last:     time.Now().UnixMilli(),
	}

	if len(start) > 0 {
		winMeta.Start = start[0]
	}

	s.windows[key][id] = window[T]{
		id:       winMeta.Id,
		metadata: winMeta.Metadata,
		elems:    []T{elem},
		start:    winMeta.Start,
		end:      winMeta.End,
		last:     winMeta.Last,
	}

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

	for id := range s.windows[k] {
		wm := WindowMeta{
			Id:       s.windows[k][id].id,
			Metadata: s.windows[k][id].metadata,
			Start:    s.windows[k][id].start,
			End:      s.windows[k][id].end,
			Last:     s.windows[k][id].last,
		}

		meta = append(meta, wm)
	}

	return meta
}

func (s *MemoryStorage[T]) GetWindowMetadata(k string, id string) WindowMeta {
	if winById, ok := s.windows[k]; ok {
		if win, ok := winById[id]; ok {
			return WindowMeta{
				Id:       win.id,
				Last:     win.last,
				End:      win.end,
				Start:    win.start,
				Metadata: win.metadata,
			}
		}
	}

	return WindowMeta{}
}

func (s *MemoryStorage[T]) CloseKeyWindows(k string, id string) {
	if winById, ok := s.windows[k]; ok {
		if win, ok := winById[id]; ok {
			win.end = time.Now().UnixMilli()
		}
	}
}

func (s *MemoryStorage[T]) CloseWindow(k string, id string) {
	if winById, ok := s.windows[k]; ok {
		if win, ok := winById[id]; ok {
			win.end = time.Now().UnixMilli()
			s.windows[k][id] = win
		}
	}
}

func (s *MemoryStorage[T]) FlushWindow(k string, id string) []T {
	if s.windows[k] == nil {
		return []T{}
	}

	if winById, ok := s.windows[k]; ok {
		if win, ok := winById[id]; ok {
			delete(s.windows[k], id)
			return win.elems
		}
	}

	return []T{}
}

func (s *MemoryStorage[T]) PushItemToWindow(k string, id string, item T) int {
	size := 0

	if winById, ok := s.windows[k]; ok {
		if win, ok := winById[id]; ok {
			win.elems = append(win.elems, item)
			s.windows[k][id] = win
			size = len(win.elems)
		}
	}

	return size
}

func (s *MemoryStorage[T]) DestroyWindow(k string, id string) {
	if winById, ok := s.windows[k]; ok {
		delete(winById, id)
	}
}

func (s *MemoryStorage[T]) StoreCheckpoint(mId string, cId string, msg T) {
	s.checkpoints[mId] = Checkpoint[T]{
		id:      cId,
		message: msg,
	}
}

func (s *MemoryStorage[T]) ClearCheckpoint(mId string) {
	delete(s.checkpoints, mId)
}

func (s *MemoryStorage[T]) GetCheckpointMessages(cId string) []T {
	msgs := []T{}

	for k := range s.checkpoints {
		if s.checkpoints[k].id == cId {
			msgs = append(msgs, s.checkpoints[k].message)
		}
	}

	return msgs
}
