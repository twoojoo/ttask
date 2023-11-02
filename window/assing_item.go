package window

import (
	"time"

	"github.com/twoojoo/ttask/storage"
	"github.com/twoojoo/ttask/task"
)

func assignItemToWindows[T any](meta []storage.WindowMeta, watermark time.Duration, item *task.Message[T]) []storage.WindowMeta {
	m := []storage.WindowMeta{}
	et := item.EventTime.UnixMilli()
	wm := watermark.Milliseconds()

	for i := range meta {
		if meta[i].Start < et && meta[i].End >= et+wm {
			m = append(m, meta[i])
		}
	}

	return m
}
