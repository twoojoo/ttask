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
//   - MaxInactivity: 1 second
//   - MaxSize: MaxInactivity * 2
type SWOptions[T any] struct {
	Id            string
	Storage       storage.Storage[task.Message[T]]
	MaxSize       time.Duration
	MaxInactivity time.Duration
}

func parseSWOptions[T any](o *SWOptions[T]) {
	if o.Storage == nil {
		o.Storage = storage.Memory[T]()
	}

	if o.Id == "" {
		o.Id = uuid.New().String()
	}

	if o.MaxInactivity == 0 {
		o.MaxInactivity = 1 * time.Second
	}

	if o.MaxSize == 0 {
		o.MaxSize = 2 * o.MaxInactivity
	}
}

func SessionWindow[T any](options SWOptions[T]) task.Operator[T, []T] {
	parseSWOptions(&options)

	return func(m *task.Meta, x *task.Message[T], next *task.Step) {
		meta := options.Storage.GetWindowsMetadata(x.Key)
		meta = filterClosedWindowMeta(meta)

		if len(meta) > 0 { // window exists
			options.Storage.PushItemToWindow(x.Key, meta[0].Id, *x)

			go func() {
				time.Sleep(options.MaxInactivity)
				meta := options.Storage.GetWindowsMetadata(x.Key)
				meta = filterClosedWindowMeta(meta)

				if meta[0].End == 0 && meta[0].Last <= time.Now().UnixMilli()-options.MaxInactivity.Milliseconds() {
					options.Storage.CloseWindow(x.Key, meta[0].Id)
					items := options.Storage.FlushWindow(x.Key, meta[0].Id)
					if len(items) > 0 {
						m.ExecNext(task.ToArray(x, items), next)
					}
				}
			}()
		} else { // window doesn't exist
			options.Storage.StartNewWindow(x.Key, *x)

			go func() {
				time.Sleep(options.MaxSize)
				meta := options.Storage.GetWindowsMetadata(x.Key)
				meta = filterClosedWindowMeta(meta)

				if meta[0].End == 0 {
					options.Storage.CloseWindow(x.Key, meta[0].Id)
					items := options.Storage.FlushWindow(x.Key, meta[0].Id)
					if len(items) > 0 {
						m.ExecNext(task.ToArray(x, items), next)
					}
				}
			}()
		}
	}
}
