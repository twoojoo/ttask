package operator

import (
	"fmt"
	"time"

	"github.com/twoojoo/ttask/task"
)

// Set a custom message key from the message itself.
func WithCustomKey[T any](extractor func(x T) string) task.Operator[T, T] {
	return func(m *task.Inner, x *task.Message[T], step *task.Step) {
		x.Key = extractor(x.Value)

		m.ExecNext(x, step)
	}
}

// Set a custom message event time from the message itself.
func WithEventTime[T any](extractor func(x T) time.Time) task.Operator[T, T] {
	return func(m *task.Inner, x *task.Message[T], step *task.Step) {
		x.EventTime = extractor(x.Value)

		m.ExecNext(x, step)
	}
}

// Print message value with a given prefix.
func Print[T any](prefix ...string) task.Operator[T, T] {
	return func(m *task.Inner, x *task.Message[T], step *task.Step) {
		if len(prefix) > 0 {
			fmt.Println(prefix[0], x.Value)
		} else {
			fmt.Println(x.Value)
		}

		m.ExecNext(x, step)
	}
}

// Print message metadata and value with a given prefix.
func PrintRaw[T any](prefix ...string) task.Operator[T, T] {
	return func(m *task.Inner, x *task.Message[T], step *task.Step) {
		if len(prefix) > 0 {
			fmt.Printf("%s %+v\n", prefix[0], x)
		} else {
			fmt.Printf("%+v\n", x)
		}

		m.ExecNext(x, step)
	}
}

// Map the message value.
func Map[T, R any](cb func(x T) R) task.Operator[T, R] {
	return func(m *task.Inner, x *task.Message[T], next *task.Step) {
		m.ExecNext(task.ReplaceValue(x, cb(x.Value)), next)
	}
}

// Map the message value (with access to task metadata and message metadata).
// Also allows to create custom operators.
func MapRaw[T, R any](cb func(m *task.Inner, x *task.Message[T]) R) task.Operator[T, R] {
	return func(m *task.Inner, x *task.Message[T], next *task.Step) {
		m.ExecNext(task.ReplaceValue(x, cb(m, x)), next)
	}
}

// Filter messages.
func Filter[T, R any](cb func(x T) bool) task.Operator[T, T] {
	return func(m *task.Inner, x *task.Message[T], next *task.Step) {
		ok := cb(x.Value)
		if ok {
			m.ExecNext(x, next)
		}
	}
}

// Filter messages (with access to task metadata and message metadata).
func FilterRaw[T, R any](cb func(m *task.Inner, x *task.Message[T]) bool) task.Operator[T, T] {
	return func(m *task.Inner, x *task.Message[T], next *task.Step) {
		ok := cb(m, x)
		if ok {
			m.ExecNext(x, next)
		}
	}
}

// Perform an action for the message.
func Tap[T any](cb func(x T)) task.Operator[T, T] {
	return func(m *task.Inner, x *task.Message[T], next *task.Step) {
		cb(x.Value)
		m.ExecNext(x, next)
	}
}

// Perform an action for the message (with access to task metadata and message metadata).
func TapRaw[T any](cb func(m *task.Inner, x *task.Message[T])) task.Operator[T, T] {
	return func(m *task.Inner, x *task.Message[T], next *task.Step) {
		cb(m, x)
		m.ExecNext(x, next)
	}
}
