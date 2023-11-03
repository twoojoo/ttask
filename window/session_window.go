package window

import (
	"time"

	"github.com/twoojoo/ttask/storage"
	"github.com/twoojoo/ttask/task"
)

func SessionWindow[T any](options SWOptions[T]) task.Operator[T, []T] {
	parseSWOptions(&options)
	storage := storage.NewStorageInterface(&options.Storage)

	//store inactivity check goroutines
	stopIncactivityCheckCh := map[string]chan int{}

	return func(inner *task.Inner, x *task.Message[T], next *task.Step) {
		meta := storage.GetWindowsMetadata(x.Key)
		mt := getMessageTime(options.WindowingTime, x)
		meta = assignMessageToWindows(meta, x, mt)

		//cancel last inactivity check for this key
		if stopIncactivityCheckCh[x.Key] != nil {
			stopIncactivityCheckCh[x.Key] <- 1
		}

		if len(meta) > 0 { // window exists
			storage.PushItemToWindow(x.Key, meta[0].Id, *x)

			//start inactivity check and store stopping channel
			stopIncactivityCheckCh[x.Key] = startInactivityCheck(options.MaxInactivity, func() {
				meta := storage.GetWindowMetadata(x.Key, meta[0].Id)

				//on incactivity: close
				if meta.End == 0 && meta.Last <= time.Now().UnixMilli()-options.MaxInactivity.Milliseconds() {
					storage.CloseWindow(x.Key, meta.Id, options.Watermark, func(items []task.Message[T]) {
						if len(items) > 0 {
							inner.ExecNext(task.ToArray(x, items), next)
						}
					})
				}
			})
		} else { // window doesn't exist
			meta := storage.StartNewWindow(x.Key, *x)

			//start inactivity check and store stopping channel
			stopIncactivityCheckCh[x.Key] = startInactivityCheck(options.MaxInactivity, func() {
				meta := storage.GetWindowMetadata(x.Key, meta.Id)

				//on incactivity: close
				if meta.End == 0 && meta.Last <= time.Now().UnixMilli()-options.MaxInactivity.Milliseconds() {
					storage.CloseWindow(x.Key, meta.Id, options.Watermark, func(items []task.Message[T]) {
						if len(items) > 0 {
							inner.ExecNext(task.ToArray(x, items), next)
						}
					})
				}
			})

			// start max size counter
			go func() {
				time.Sleep(options.MaxSize)

				meta := storage.GetWindowMetadata(x.Key, meta.Id)

				//on max size: close
				if meta.End == 0 {
					storage.CloseWindow(x.Key, meta.Id, options.Watermark, func(items []task.Message[T]) {
						if len(items) > 0 {
							inner.ExecNext(task.ToArray(x, items), next)
						}
					})
				}
			}()
		}
	}
}
