package ttask

import (
	"log"
	"time"
)

func assignMessageToWindows[T any](
	meta []windowMeta,
	item *Message[T],
	messageTime time.Time,
) []windowMeta {
	m := []windowMeta{}

	for i := range meta {
		afterStart := messageTime.After(meta[i].Start) || meta[i].Start.Equal(meta[i].Start)
		beforeEnd := meta[i].End.IsZero() || messageTime.Before(meta[i].End)

		if afterStart && beforeEnd {
			m = append(m, meta[i])
		}
	}

	return m
}

func getMessageTime[T any](wTime WindowingTime, msg *Message[T]) time.Time {
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
