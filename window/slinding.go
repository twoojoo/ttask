package window

import (
	"time"

	"github.com/twoojoo/ttask/storage"
	"github.com/twoojoo/ttask/task"
)

type SWOptions[T any] struct {
	Id      string
	Storage storage.Storage[task.Message[T]]
	Size    time.Duration
}

// SlidingWindow:
//
//..0....1....2....3.........4.........5....6....7...
//
//[-------------][-------------][-------------][-----
func SlidingWindow[T any](options SWOptions[T]) task.Operator[T, []T] {
	first := true

	return func(m *task.Meta, x *task.Message[T], next *task.Step) {
		if first {
			recovery(options, func(key string) {
				flush(options, key, m, x, next)
			})

			go startTicker[T](options, func(tick time.Time) {
				keys := options.Storage.GetAllKeys()

				now := time.Now()
				start := now.UnixMicro()
				end := now.Add(options.Size).UnixMicro()

				for _, k := range keys {
					flush(options, k, m, x, next)

					options.Storage.SetMetadata(k, map[string]int64{
						"start": start,
						"end":   end,
					})
				}
			})

			first = false
		}

		md := options.Storage.GetMetadata(x.Key)
		newMd := md

		now := time.Now()
		start := now.UnixMicro()
		end := now.Add(options.Size).UnixMicro()

		if newMd["start"] == 0 || newMd["end"] == 0 {
			newMd["start"] = start
			newMd["end"] = end
		}

		(options.Storage).Push(x.Key, x, newMd)
	}
}

func startTicker[T any](options SWOptions[T], onTick func(tick time.Time)) {
	t := time.NewTicker(options.Size)
	defer t.Stop()

	for tick := range t.C {
		onTick(tick)
	}
}

//check if there are some open windows that should be ended and flush them
func recovery[T any](options SWOptions[T], flush func(key string)) {
	now := time.Now()

	sizes := options.Storage.GetAllSizes()

	for key, size := range sizes {
		if size > 0 {
			md := options.Storage.GetMetadata(key)

			if md["end"] >= now.UnixMicro() {
				flush(key)
			}
		}
	}
}

func flush[T any](options SWOptions[T], k string, m *task.Meta, x *task.Message[T], next *task.Step) {
	items := options.Storage.Flush(k)
	if len(items) > 0 {
		m.ExecNext(task.ToArray(x, items), next)
	}
}
