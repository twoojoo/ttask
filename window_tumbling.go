package ttask

import (
	"time"
)

// TumblingWindow:
//
// ..0....1....2....3.........4.........5....6....7...
//
// [-------------][-------------][-------------][-----
func TumblingWindow[T any](id string, options TWOptions[T]) Operator[T, []T] {
	parseTWOptions(&options)

	first := true

	return func(i *Inner, x *Message[T], next *Step) {
		sw := newStorageWrapper[T](i)

		if first {
			go func() {
				i.wg.Add(1)
				for range time.Tick(options.Size) {

					now := time.Now()

					keys, err := sw.getKeys(id)
					if err != nil {
						i.Error(err)
						return
					}

					for _, k := range keys {
						meta, err := sw.getWindowsMetadataByKey(id, k)
						if err != nil {
							i.Error(err)
							return
						}

						for j := range meta {
							sw.closeWindow(id, x.Key, meta[j].Id, options.Watermark, func(items []Message[T]) {
								if len(items) > 0 {
									i.ExecNext(toArray(x, items), next)
								}
							})
						}

						sw.startNewEmptyWindow(id, k, now)
					}

					first = false

					// i.wg.Done()
				}
			}()
		}

		for {
			if !first {
				break
			}
		}

		meta, err := sw.getWindowsMetadataByKey(id, x.Key)
		if err != nil {
			i.Error(err)
			return
		}

		mt := getMessageTime(options.WindowingTime, x)
		meta = assignMessageToWindows(meta, x, mt)

		if len(meta) == 0 {
			_, err := sw.startNewWindow(id, x.Key, *x)
			if err != nil {
				i.Error(err)
				return
			}
		} else {
			_, err := sw.pushMessageToWindow(id, x.Key, meta[0].Id, *x)
			if err != nil {
				i.Error(err)
				return
			}
		}
	}
}
