package window

import (
	"time"
)

func startInactivityCheck(maxInactivity time.Duration, onInactive func()) chan int {
	//2 because the first message won't block the calling routine
	ch := make(chan int, 2)

	go func() {
		time.Sleep(maxInactivity)

		select {
		case <- ch:
			return
		default:
			onInactive()
		}
	}()

	return ch
}
