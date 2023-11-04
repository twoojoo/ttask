package ttask

import (
	"time"
)

func HoppingWindow[T any](id string, options HWOptions[T]) Operator[T, []T] {
	parseHWOptions(&options)

	first := true
	var nextStart time.Time

	return func(i *Inner, x *Message[T], next *Step) {
		if first {
			first = false
			go startWinLoop[T](options, func() {
				start := time.Now()
				nextStart = start

				keys, err := getKeys(i.storage, id)
				if err != nil {
					i.Error(err)
					return
				}

				for _, k := range keys {
					_, err := startNewEmptyWindow(i.storage, id, k, nextStart)
					if err != nil {
						i.Error(err)
						return
					}
				}

			}, func() {
				end := time.Now()

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
						if meta[j].End.IsZero() && (meta[j].Start.Before(end.Add(-options.Size)) || meta[j].Start.Equal(end.Add(-options.Size))) {
							closeWindow(i.storage, id, x.Key, meta[j].Id, options.Watermark, func(items []Message[T]) {
								if len(items) > 0 {
									i.ExecNext(toArray(x, items), next)
								}
							})
						}
					}
				}
			})
		}

		//pushing item

		//wait for nextStart to be set by the loop
		for {
			if !nextStart.IsZero() {
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

		// if no window for this key, just create 1 with the last start ts
		if len(meta) == 0 && !first {
			_, err := startNewWindow(i.storage, id, x.Key, *x, nextStart)
			if err != nil {
				i.Error(err)
				return
			}
		} else {
			lastExists := false

			// push item to all windows for that key that are not closed yet
			for _, m := range meta {
				if m.End.IsZero() {
					_, err := pushMessageToWindow(i.storage, id, x.Key, m.Id, *x)
					if err != nil {
						i.Error(err)
						return
					}
				}

				// check if the next window is already created
				if m.Start.After(nextStart) || m.Start.Equal(nextStart) {
					lastExists = true
				}
			}

			// if next window is not yet created, then create it
			if !lastExists {
				_, err := startNewWindow(i.storage, id, x.Key, *x, nextStart)
				if err != nil {
					i.Error(err)
					return
				}

			}
		}
	}
}

func startWinLoop[T any](options HWOptions[T], onStart func(), onClose func()) {
	for {
		onStart()

		go func() {
			time.Sleep(options.Size)
			onClose()
		}()

		time.Sleep(options.Hop)
	}
}
