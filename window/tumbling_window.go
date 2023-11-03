package window

import (
	"time"

	"github.com/twoojoo/ttask/storage"
	"github.com/twoojoo/ttask/task"
)

// TumblingWindow:
//
// ..0....1....2....3.........4.........5....6....7...
//
// [-------------][-------------][-------------][-----
func TumblingWindow[T any](options TWOptions[T]) task.Operator[T, []T] {
	parseTWOptions(&options)

	storage := storage.NewStorageInterface(&options.Storage)

	first := true

	return func(inner *task.Inner, x *task.Message[T], next *task.Step) {
		if first {
			go func() {
				for range time.Tick(options.Size) {
					now := time.Now().UnixMilli()

					keys := storage.GetKeys()

					for _, k := range keys {
						meta := storage.GetWindowsMetadata(k)

						for i := range meta {
							storage.CloseWindow(x.Key, meta[i].Id, options.Watermark, func(items []task.Message[T]) {
								if len(items) > 0 {
									inner.ExecNext(task.ToArray(x, items), next)
								}
							})
						}

						storage.StartNewEmptyWindow(k, now)
					}

					first = false
				}
			}()
		}

		for {
			if !first {
				break
			}
		}

		meta := storage.GetWindowsMetadata(x.Key)
		mt := getMessageTime(options.WindowingTime, x)
		meta = assignMessageToWindows(meta, x, mt)

		if len(meta) == 0 {
			storage.StartNewWindow(x.Key, *x)
		} else {
			storage.PushItemToWindow(x.Key, meta[0].Id, *x)
		}
	}
}

func filterClosedWindowMeta(meta []storage.WindowMeta) []storage.WindowMeta {
	filtered := []storage.WindowMeta{}

	for i := range meta {
		if meta[i].End == 0 {
			filtered = append(filtered, meta[i])
		}
	}

	return filtered
}
