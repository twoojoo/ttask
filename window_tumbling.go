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
		if first {
			go func() {
				for range time.Tick(options.Size) {
					now := time.Now()

					keys, err := getKeys(i.storage, id)
					if err != nil {
						i.Error(err)
						return
					}

					for _, k := range keys {
						meta, err := getWindowsMetadataByKey(i.storage, id, k)
						if err != nil {
							i.Error(err)
							return
						}

						for j := range meta {
							closeWindow(i.storage, id, x.Key, meta[j].Id, options.Watermark, func(items []Message[T]) {
								if len(items) > 0 {
									i.ExecNext(toArray(x, items), next)
								}
							})
						}

						startNewEmptyWindow(i.storage, id, k, now)
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

		meta, err := getWindowsMetadataByKey(i.storage, id, x.Key)
		if err != nil {
			i.Error(err)
			return
		}

		mt := getMessageTime(options.WindowingTime, x)
		meta = assignMessageToWindows(meta, x, mt)

		if len(meta) == 0 {
			_, err := startNewWindow(i.storage, id, x.Key, *x)
			if err != nil {
				i.Error(err)
				return
			}
		} else {
			_, err := pushMessageToWindow(i.storage, id, x.Key, meta[0].Id, *x)
			if err != nil {
				i.Error(err)
				return
			}
		}
	}
}

func filterClosedWindowMeta(meta []windowMeta) []windowMeta {
	filtered := []windowMeta{}

	for i := range meta {
		if meta[i].End.IsZero() {
			filtered = append(filtered, meta[i])
		}
	}

	return filtered
}
