package window

import (
	"log"
	"time"

	"github.com/twoojoo/ttask/storage"
	"github.com/twoojoo/ttask/task"
)

func assignMessageToWindows[T any](
	meta []storage.WindowMeta,
	item *task.Message[T],
	messageTime time.Time,
) []storage.WindowMeta {
	m := []storage.WindowMeta{}
	t := messageTime.UnixMilli()

	for i := range meta {
		if t >= meta[i].Start && (meta[i].End == 0 || t < meta[i].End) {
			m = append(m, meta[i])
		}
	}

	return m
}

func getMessageTime[T any](wTime WindowingTime, msg *task.Message[T]) time.Time {
	switch wTime {
	case EventTime:
		if msg.EventTime.IsZero() {
			log.Println("!> event time not set - fallback on processing time")
			return time.Now()
		}
	case InjestionTime:
		return msg.GetInjestionTime()
	}

	return time.Now()
}
