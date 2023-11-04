package ttask

import (
	"time"

	"github.com/google/uuid"
)

type MemoryStorage struct {
	windows map[string]map[string]map[string]window
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		windows: map[string]map[string]map[string]window{},
	}
}

func (s *MemoryStorage) storeCheckpoint(taskId string, msgId string, cpId string, msg []byte) error {
	return nil
}

func (s *MemoryStorage) clearCheckpoint(taskId string, msgId string, cpId string) error {
	return nil
}

func (s *MemoryStorage) getCheckpointMessages(taskId string, cpId string) ([][]byte, error) {
	return [][]byte{}, nil
}

func (s *MemoryStorage) startNewEmptyWindow(cpId string, key string, start ...time.Time) (windowMeta, error) {
	if s.windows[cpId][key] == nil {
		s.windows[cpId][key] = make(map[string]window, 0)
	}

	id := uuid.NewString()

	winMeta := windowMeta{
		Id:    id,
		Start: time.Now(),
	}

	if len(start) > 0 {
		winMeta.Start = start[0]
	}

	s.windows[cpId][key][id] = window{
		id:    winMeta.Id,
		elems: []any{},
		start: winMeta.Start,
		end:   winMeta.End,
	}

	return winMeta, nil
}

func (s *MemoryStorage) startNewWindow(cpId string, key string, msg any, start ...time.Time) (windowMeta, error) {
	if s.windows[cpId] == nil {
		s.windows[cpId] = map[string]map[string]window{}
	}
	if s.windows[cpId][key] == nil {
		s.windows[cpId][key] = map[string]window{}
	}

	id := uuid.NewString()

	winMeta := windowMeta{
		Id:    id,
		Start: time.Now(),
		Last:  time.Now(),
	}

	if len(start) > 0 {
		winMeta.Start = start[0]
	}

	s.windows[cpId][key][id] = window{
		id:    winMeta.Id,
		elems: []any{msg},
		start: winMeta.Start,
		end:   winMeta.End,
		last:  winMeta.Last,
	}

	return winMeta, nil
}

func (s *MemoryStorage) pushMessageToWindow(cpId string, key string, winId string, msg any) (int, error) {
	size := 0

	if winByCpId, ok := s.windows[cpId]; ok {
		if winById, ok := winByCpId[key]; ok {
			if win, ok := winById[winId]; ok {
				win.elems = append(win.elems, msg)
				s.windows[cpId][key][winId] = win
				size = len(win.elems)
			}
		}
	}

	return size, nil
}

func (s *MemoryStorage) getWindowMetadata(cpId string, key string, winId string) (windowMeta, error) {
	if winByCpId, ok := s.windows[cpId]; ok {
		if winById, ok := winByCpId[key]; ok {
			if win, ok := winById[winId]; ok {
				return windowMeta{
					Id:    win.id,
					Last:  win.last,
					End:   win.end,
					Start: win.start,
				}, nil
			}
		}
	}

	return windowMeta{}, nil
}

func (s *MemoryStorage) getWindowsMetadataByKey(cpId string, key string) ([]windowMeta, error) {
	if s.windows[cpId] == nil {
		return []windowMeta{}, nil
	}

	if s.windows[cpId][key] == nil {
		return []windowMeta{}, nil
	}

	meta := []windowMeta{}

	for id := range s.windows[cpId][key] {
		wm := windowMeta{
			Id:    s.windows[cpId][key][id].id,
			Start: s.windows[cpId][key][id].start,
			End:   s.windows[cpId][key][id].end,
			Last:  s.windows[cpId][key][id].last,
		}

		meta = append(meta, wm)
	}

	return meta, nil
}

func (s *MemoryStorage) closeWindow(cpId string, key string, winId string) error {
	if winByCpId, ok := s.windows[cpId]; ok {
		if winById, ok := winByCpId[key]; ok {
			if win, ok := winById[winId]; ok {
				win.end = time.Now()
				s.windows[cpId][key][winId] = win
			}
		}
	}

	return nil
}

func (s *MemoryStorage) flushWindow(cpId string, key string, winId string) ([]any, error) {
	if s.windows[cpId] == nil {
		return []any{}, nil
	}

	if s.windows[cpId][key] == nil {
		return []any{}, nil
	}

	if winByCpId, ok := s.windows[cpId]; ok {
		if winById, ok := winByCpId[key]; ok {
			if win, ok := winById[winId]; ok {
				delete(s.windows[cpId][key], winId)
				return win.elems, nil
			}
		}
	}

	return []any{}, nil
}

func (s *MemoryStorage) destroyWindow(cpId string, key string, winId string) error {
	if winById, ok := s.windows[key]; ok {
		delete(winById, winId)
	}

	return nil
}

func (s *MemoryStorage) getWindowSize(cpId string, key string, winId string) (int, error) {
	if s.windows[cpId] == nil {
		return 0, nil
	}

	if s.windows[cpId][key] == nil {
		return 0, nil
	}

	for _, win := range s.windows[cpId][key] {
		if win.id == winId {
			return len(win.elems), nil
		}
	}

	return 0, nil
}

func (s *MemoryStorage) getKeys(cpId string) ([]string, error) {
	keys := []string{}

	if s.windows[cpId] == nil {
		return keys, nil
	}

	for k := range s.windows[cpId] {
		keys = append(keys, k)
	}

	return keys, nil
}
