package window

import (
	"time"

	"github.com/google/uuid"
	"github.com/twoojoo/ttask/storage"
	"github.com/twoojoo/ttask/task"
)

// Defaults:
//   - Storage: memory (no persistence)
//   - Id: random uuid
//   - Size: 1 second
//   - Hop: 2 seconds
type TWOptions[T any] struct {
	Id        string
	Storage   storage.Storage[task.Message[T]]
	Size      time.Duration
	Watermark time.Duration
}

func parseTWOptions[T any](o *TWOptions[T]) {
	if o.Storage == nil {
		o.Storage = storage.Memory[T]()
	}

	if o.Id == "" {
		o.Id = uuid.New().String()
	}

	if o.Size == 0 {
		o.Size = 1 * time.Second
	}
}

// TumblingWindow:
//
// ..0....1....2....3.........4.........5....6....7...
//
// [-------------][-------------][-------------][-----
func TumblingWindow[T any](options TWOptions[T]) task.Operator[T, []T] {
	parseTWOptions(&options)
	storage := storage.NewStorageInterface(&options.Storage)

	first := true

	return func(m *task.Meta, x *task.Message[T], next *task.Step) {
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
									m.ExecNext(task.ToArray(x, items), next)
								}
							})

							// idsToFlush = append(idsToFlush, meta[i].Id)
						}

						storage.StartNewEmptyWindow(k, now)
					}

					first = false

					// for _, k := range keys {
					// 	for _, id := range idsToFlush {
					// 		items := storage.FlushWindow(k, id)
					// 		if len(items) > 0 {
					// 			m.ExecNext(task.ToArray(x, items), next)
					// 		}
					// 	}
					// }
				}
			}()
		}

		for {
			if !first {
				break
			}
		}

		meta := storage.GetWindowsMetadata(x.Key)
		meta = filterClosedWindowMeta(meta)

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
