package ttask

// Counting Window:
//
// ...1....2.........3...........4...5......6........7....8....
//
// ..[----------------].........[------------]......[----------
func CountingWindow[T any](id string, options CWOptions[T]) Operator[T, []T] {
	parseCWOptions(&options)

	stopIncactivityCheckCh := map[string]chan int{}

	return func(i *Inner, x *Message[T], next *Step) {
		sw := newStorageWrapper[T](i)

		//cancel last inactivity check
		if stopIncactivityCheckCh[x.Key] != nil {
			stopIncactivityCheckCh[x.Key] <- 1
		}

		meta, err := sw.getWindowsMetadataByKey(id, x.Key)
		if err != nil {
			i.Error(err)
			return
		}

		var size int
		if len(meta) > 1 {
			panic("there should be only 1 window per key in counting window")
		} else if len(meta) == 0 {
			newWinMeta, err := sw.startNewWindow(id, x.Key, *x)
			if err != nil {
				i.Error(err)
				return
			}

			meta = append(meta, newWinMeta)
			size = 1
		} else {
			size, err = sw.pushMessageToWindow(id, x.Key, meta[0].Id, *x)
			if err != nil {
				i.Error(err)
				return
			}
		}

		// start new inactivity check
		if options.MaxInactivity > 0 && options.Size > 1 {
			stopIncactivityCheckCh[x.Key] = startInactivityCheck(i, options.MaxInactivity, func() {
				items, err := sw.flushWindow(id, x.Key, meta[0].Id)
				if err != nil {
					i.Error(err)
					return
				}

				if len(items) > 0 {
					i.ExecNext(toArray(x, items), next)
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
			items, err := sw.flushWindow(id, x.Key, meta[0].Id)
			if err != nil {
				i.Error(err)
				return
			}

			if len(items) > 0 {
				i.ExecNext(toArray(x, items), next)
			}
		}
	}
}
