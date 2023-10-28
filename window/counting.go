package window

import (
	"time"

	"github.com/twoojoo/ttask/storage"
	"github.com/twoojoo/ttask/task"
)

type CWOptions[T any] struct {
	Id            string
	Storage       storage.Storage[task.Message[T]]
	Size          int
	MaxInactivity time.Duration
}

// Counting Window:
//
// ...1....2.........3...........4...5......6........7....8....
//
// ..[----------------].........[------------]......[----------
func CountingWindow[T any](options CWOptions[T]) task.Operator[T, []T] {
	stopIncactivityCheckCh := map[string]chan int{}

	return func(m *task.Meta, x *task.Message[T], next *task.Step) {
		//cancel last inactivity check
		if stopIncactivityCheckCh[x.Key] != nil {
			stopIncactivityCheckCh[x.Key] <- 1
		}

		size := options.Storage.GetSize(x.Key)

		md := map[string]int64{}
		if size == 0 {
			now := time.Now()
			start := now.UnixMicro()
			md = map[string]int64{"start": start}
		}

		size = (options.Storage).Push(x.Key, x, md)

		// start new inactivity check
		if options.MaxInactivity > 0 && options.Size > 1 {
			stopIncactivityCheckCh[x.Key] = startInactivityCheck(options.MaxInactivity, func() {
				items := (options.Storage).Flush(x.Key)
				if len(items) > 0 {
					m.ExecNext(task.ToArray(x, items), next)
				}
			})
		}

		// normal flush
		if size >= options.Size && options.Size != 0 {
			//cancel last inactivity check
			if stopIncactivityCheckCh[x.Key] != nil {
				stopIncactivityCheckCh[x.Key] <- 1
			}

			items := (options.Storage).Flush(x.Key)
			if len(items) > 0 {
				m.ExecNext(task.ToArray(x, items), next)
			}
		}
	}
}
