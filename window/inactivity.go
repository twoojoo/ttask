package window

import (
	"time"
)

func startInactivityCheck(maxInactivity time.Duration, onInactive func()) chan int {
	ch := make(chan int, 1)

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
