package window

import (
	"sync"
	"time"

	"github.com/twoojoo/ttask/storage"
	"github.com/twoojoo/ttask/task"
)

type TWCOptions[T any] struct {
	Id            string
	Storage       storage.Storage[task.Message[T]]
	Size          int
	MaxInactivity time.Duration
}

func TumblingWindowCount[T any](options TWCOptions[T]) task.Operator[T, []T] {
	var mu sync.Mutex

	return func(m *task.Meta, x *task.Message[T], next *task.Step) {
		size := (options.Storage).Push(x.Key, x)

		if options.MaxInactivity > 0 && size == 1 && options.Size > 1 {
			go func() {
				time.Sleep(options.MaxInactivity)

				mu.Lock()

				items := (options.Storage).Flush(x.Key)
				if len(items) > 0 {
					m.ExecNext(task.ToArray(x, items), next)
				}

				mu.Unlock()
			}()
		}

		if size >= options.Size && options.Size != 0 {
			mu.Lock()

			items := (options.Storage).Flush(x.Key)

			if len(items) > 0 {
				m.ExecNext(task.ToArray(x, items), next)
			}

			mu.Unlock()
		}
	}
}
