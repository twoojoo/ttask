package ttask

import (
	"time"
)

func SessionWindow[T any](id string, options SWOptions[T]) Operator[T, []T] {
	parseSWOptions(&options)

	//store inactivity check goroutines
	stopIncactivityCheckCh := map[string]chan int{}

	return func(i *Inner, x *Message[T], next *Step) {
		sw := newStorageWrapper[T](i)

		meta, err := sw.getWindowsMetadataByKey(id, x.Key)
		if err != nil {
			i.Error(err)
			return
		}

		mt := getMessageTime(options.WindowingTime, x)
		meta = assignMessageToWindows(meta, x, mt)

		//cancel last inactivity check for this key
		if stopIncactivityCheckCh[x.Key] != nil {
			stopIncactivityCheckCh[x.Key] <- 1
		}

		if len(meta) > 0 { // window exists
			_, err := sw.pushMessageToWindow(id, x.Key, meta[0].Id, *x)
			if err != nil {
				i.Error(err)
				return
			}

			//start inactivity check and store stopping channel
			stopIncactivityCheckCh[x.Key] = startInactivityCheck(i, options.MaxInactivity, func() {
				meta, err := sw.getWindowMetadata(id, x.Key, meta[0].Id)
				if err != nil {
					i.Error(err)
					return
				}

				//on incactivity: close
				if meta.End.IsZero() && (meta.Last.Before(time.Now().Add(-options.MaxInactivity)) || meta.Last.Equal(time.Now().Add(-options.MaxInactivity))) {
					err := sw.closeWindow(id, x.Key, meta.Id, options.Watermark, func(items []Message[T]) {
						if len(items) > 0 {
							i.ExecNext(toArray(x, items), next)
						}
					})

					if err != nil {
						i.Error(err)
						return
					}
				}
			})
		} else { // window doesn't exist
			meta, err := sw.startNewWindow(id, x.Key, *x)
			if err != nil {
				i.Error(err)
				return
			}

			//start inactivity check and store stopping channel
			stopIncactivityCheckCh[x.Key] = startInactivityCheck(i, options.MaxInactivity, func() {
				meta, err := sw.getWindowMetadata(id, x.Key, meta.Id)
				if err != nil {
					i.Error(err)
					return
				}

				//on incactivity: close
				if meta.End.IsZero() && (meta.Last.Before(time.Now().Add(-options.MaxInactivity)) || meta.Last.Equal(time.Now().Add(-options.MaxInactivity))) {
					err := sw.closeWindow(id, x.Key, meta.Id, options.Watermark, func(items []Message[T]) {
						if len(items) > 0 {
							i.ExecNext(toArray(x, items), next)
						}
					})

					if err != nil {
						i.Error(err)
						return
					}
				}
			})

			// start max size counter
			go func() {
				time.Sleep(options.MaxSize)

				meta, err := sw.getWindowMetadata(id, x.Key, meta.Id)
				if err != nil {
					i.Error(err)
					return
				}

				//on max size: close
				if meta.End.IsZero() {
					err := sw.closeWindow(id, x.Key, meta.Id, options.Watermark, func(items []Message[T]) {
						if len(items) > 0 {
							i.ExecNext(toArray(x, items), next)
						}
					})

					if err != nil {
						i.Error(err)
						return
					}
				}
			}()
		}
	}
}
