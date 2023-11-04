package ttask

import (
	"fmt"
	"time"
)

// Set a custom message key from the message itself.
func WithCustomKey[T any](extractor func(x T) string) Operator[T, T] {
	return func(inner *Inner, x *Message[T], step *Step) {
		x.Key = extractor(x.Value)

		inner.ExecNext(x, step)
	}
}

// Set a custom message event time from the message itself.
func WithEventTime[T any](extractor func(x T) time.Time) Operator[T, T] {
	return func(inner *Inner, x *Message[T], step *Step) {
		x.EventTime = extractor(x.Value)

		inner.ExecNext(x, step)
	}
}

// Print message value with a given prefix.
func Print[T any](prefix ...string) Operator[T, T] {
	return func(inner *Inner, x *Message[T], step *Step) {
		if len(prefix) > 0 {
			fmt.Println(prefix[0], x.Value)
		} else {
			fmt.Println(x.Value)
		}

		inner.ExecNext(x, step)
	}
}

// Print message metadata and value with a given prefix.
func PrintRaw[T any](prefix ...string) Operator[T, T] {
	return func(inner *Inner, x *Message[T], step *Step) {
		if len(prefix) > 0 {
			fmt.Printf("%s %+v\n", prefix[0], x)
		} else {
			fmt.Printf("%+v\n", x)
		}

		inner.ExecNext(x, step)
	}
}

// Map the message value.
func Map[T, R any](cb func(x T) R) Operator[T, R] {
	return func(inner *Inner, x *Message[T], next *Step) {
		inner.ExecNext(replaceValue(x, cb(x.Value)), next)
	}
}

// Map the message value (with access to task metadata and message metadata).
// Also allows to create custom operators.
func MapRaw[T, R any](cb func(inner *Inner, x *Message[T]) R) Operator[T, R] {
	return func(inner *Inner, x *Message[T], next *Step) {
		inner.ExecNext(replaceValue(x, cb(inner, x)), next)
	}
}

// Filter messages.
func Filter[T, R any](cb func(x T) bool) Operator[T, T] {
	return func(inner *Inner, x *Message[T], next *Step) {
		ok := cb(x.Value)
		if ok {
			inner.ExecNext(x, next)
		}
	}
}

// Filter messages (with access to task metadata and message metadata).
func FilterRaw[T, R any](cb func(inner *Inner, x *Message[T]) bool) Operator[T, T] {
	return func(inner *Inner, x *Message[T], next *Step) {
		ok := cb(inner, x)
		if ok {
			inner.ExecNext(x, next)
		}
	}
}

// Perform an action for the message.
func Tap[T any](cb func(x T)) Operator[T, T] {
	return func(inner *Inner, x *Message[T], next *Step) {
		cb(x.Value)
		inner.ExecNext(x, next)
	}
}

// Perform an action for the message (with access to task metadata and message metadata).
func TapRaw[T any](cb func(inner *Inner, x *Message[T])) Operator[T, T] {
	return func(inner *Inner, x *Message[T], next *Step) {
		cb(inner, x)
		inner.ExecNext(x, next)
	}
}
