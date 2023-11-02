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
type HWOptions[T any] struct {
	Id      string
	Storage storage.Storage[task.Message[T]]
	Size    time.Duration
	Hop     time.Duration
}

func parseHWOptions[T any](o *HWOptions[T]) {
	if o.Storage == nil {
		o.Storage = storage.Memory[T]()
	}

	if o.Id == "" {
		o.Id = uuid.New().String()
	}

	if o.Size == 0 {
		o.Size = 1 * time.Second
	}

	if o.Hop == 0 {
		o.Hop = 2 * time.Second
	}
}

func HoppingWindow[T any](options HWOptions[T]) task.Operator[T, []T] {
	parseHWOptions(&options)
	storage := storage.NewStorageInterface(&options.Storage)

	first := true
	nextStart := int64(0)

	return func(m *task.Meta, x *task.Message[T], next *task.Step) {
		if first {
			first = false
			go startWinLoop[T](options, func(start int64) {
				nextStart = start

				for _, k := range storage.GetKeys() {
					storage.StartNewEmptyWindow(k, nextStart)
				}

			}, func(end int64) {
				keys := storage.GetKeys()

				for _, k := range keys {
					meta := storage.GetWindowsMetadata(k)
					idsToFlush := []string{}

					for i := range meta {
						if meta[i].End == 0 && meta[i].Start <= (end-options.Size.Milliseconds()) {
							storage.CloseWindow(k, meta[i].Id)
							idsToFlush = append(idsToFlush, meta[i].Id)
						}
					}

					//fluhs only windows closed in this turn
					for _, id := range idsToFlush {
						items := storage.FlushWindow(k, id)
						// log.Println("fglushin window", k, id, "len", len(items))
						if len(items) > 0 {
							m.ExecNext(task.ToArray(x, items), next)
						}
					}
				}
			})
		}

		//pushing item

		//wait for nextStart to be set by the loop
		for {
			if nextStart != 0 {
				break
			}
		}

		meta := storage.GetWindowsMetadata(x.Key)

		// if no window for this key, just create 1 with the last start ts
		if len(meta) == 0 && !first {
			storage.StartNewWindow(x.Key, *x, nextStart)
		} else {
			lastExists := false

			// push item to all windows for that key that are not closed yet
			for _, m := range meta {
				if m.End == 0 {
					// log.Printf("pushing %v to %s\n", x.Value, m.Id)
					storage.PushItemToWindow(x.Key, m.Id, *x)
				}

				// check if the next window is already created
				if m.Start >= nextStart {
					lastExists = true
				}
			}

			// if next window is not yet created, then create it
			if !lastExists {
				storage.StartNewWindow(x.Key, *x, nextStart)
			}
		}
	}
}

func startWinLoop[T any](options HWOptions[T], onStart func(start int64), onClose func(end int64)) {
	for {
		onStart(time.Now().UnixMilli())

		go func() {
			time.Sleep(options.Size)
			onClose(time.Now().UnixMilli())
		}()

		time.Sleep(options.Hop)
	}
}
