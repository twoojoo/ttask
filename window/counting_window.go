package window

import (
	"github.com/twoojoo/ttask/storage"
	"github.com/twoojoo/ttask/task"
)



// Counting Window:
//
// ...1....2.........3...........4...5......6........7....8....
//
// ..[----------------].........[------------]......[----------
func CountingWindow[T any](options CWOptions[T]) task.Operator[T, []T] {
	parseCWOptions(&options)
	storage := storage.NewStorageInterface(&options.Storage)

	stopIncactivityCheckCh := map[string]chan int{}

	return func(m *task.Meta, x *task.Message[T], next *task.Step) {
		//cancel last inactivity check
		if stopIncactivityCheckCh[x.Key] != nil {
			stopIncactivityCheckCh[x.Key] <- 1
		}

		meta := storage.GetWindowsMetadata(x.Key)

		var size int
		if len(meta) > 1 {
			panic("there should be only 1 window per key in counting window")
		} else if len(meta) == 0 {
			meta = append(meta, storage.StartNewWindow(x.Key, *x))
			size = 1
		} else {
			size = storage.PushItemToWindow(x.Key, meta[0].Id, *x)
		}

		// start new inactivity check
		if options.MaxInactivity > 0 && options.Size > 1 {
			stopIncactivityCheckCh[x.Key] = startInactivityCheck(options.MaxInactivity, func() {
				// storage.CloseWindow(x.Key, meta[0].Id)
				items := storage.FlushWindow(x.Key, meta[0].Id)
				if len(items) > 0 {
					m.ExecNext(task.ToArray(x, items), next)
				}
			})
		}

		// normal flush
		if size >= options.Size && options.Size != 0 {
			//cancel last inactivity check
			if stopIncactivityCheckCh[x.Key] != nil {
				select {
				case stopIncactivityCheckCh[x.Key] <- 1:
				default:
				}
			}

			// storage.CloseWindow(x.Key, meta[0].Id)
			items := storage.FlushWindow(x.Key, meta[0].Id)
			if len(items) > 0 {
				m.ExecNext(task.ToArray(x, items), next)
			}
		}
	}
}
