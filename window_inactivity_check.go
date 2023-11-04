package ttask

import (
	"time"
)

func startInactivityCheck(inner *Inner, maxInactivity time.Duration, onInactive func()) chan int {
	//2 because the first message won't block the calling routine
	ch := make(chan int, 2)
	inner.wg.Add(1)

	go func() {
		defer func() {
			inner.wg.Done()
		}()

		time.Sleep(maxInactivity)

		select {
		case <-ch:
			return
		default:
			onInactive()
		}
	}()

	return ch
}

// func stopInactivityCheck()
