package ttask

import (
	"time"
)

func SessionWindow[T any](id string, options SWOptions[T]) Operator[T, []T] {
	parseSWOptions(&options)

	//store inactivity check goroutines
	stopIncactivityCheckCh := map[string]chan int{}

	return func(i *Inner, x *Message[T], next *Step) {
		meta, err := getWindowsMetadataByKey(i.storage, id, x.Key)
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
			_, err := pushMessageToWindow(i.storage, id, x.Key, meta[0].Id, *x)
			if err != nil {
				i.Error(err)
				return
			}

			//start inactivity check and store stopping channel
			stopIncactivityCheckCh[x.Key] = startInactivityCheck(options.MaxInactivity, func() {
				meta, err := getWindowMetadata(i.storage, id, x.Key, meta[0].Id)
				if err != nil {
					i.Error(err)
					return
				}

				//on incactivity: close
				if meta.End.IsZero() && (meta.Last.Before(time.Now().Add(-options.MaxInactivity)) || meta.Last.Equal(time.Now().Add(-options.MaxInactivity))) {
					err := closeWindow(i.storage, id, x.Key, meta.Id, options.Watermark, func(items []Message[T]) {
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
			meta, err := startNewWindow(i.storage, id, x.Key, *x)
			if err != nil {
				i.Error(err)
				return
			}

			//start inactivity check and store stopping channel
			stopIncactivityCheckCh[x.Key] = startInactivityCheck(options.MaxInactivity, func() {
				meta, err := getWindowMetadata(i.storage, id, x.Key, meta.Id)
				if err != nil {
					i.Error(err)
					return
				}

				//on incactivity: close
				if meta.End.IsZero() && (meta.Last.Before(time.Now().Add(-options.MaxInactivity)) || meta.Last.Equal(time.Now().Add(-options.MaxInactivity))) {
					err := closeWindow(i.storage, id, x.Key, meta.Id, options.Watermark, func(items []Message[T]) {
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

				meta, err := getWindowMetadata(i.storage, id, x.Key, meta.Id)
				if err != nil {
					i.Error(err)
					return
				}

				//on max size: close
				if meta.End.IsZero() {
					err := closeWindow(i.storage, id, x.Key, meta.Id, options.Watermark, func(items []Message[T]) {
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
