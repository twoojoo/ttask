package ttask

import (
	"strings"
	"time"
)

func fromItem[T any](item T) Operator[any, T] {
	return func(inner *Inner, x *Message[any], next *Step) {
		inner.ExecNext(NewMessage(item), next)
	}
}

// Source: trigger a task once with the give item.
func FromItem[T any](taskId string, item T) *TTask[any, T] {
	return Via(Task[any](taskId), fromItem(item))
}

func fromArray[T any](array []T) Operator[any, T] {
	return func(inner *Inner, x *Message[any], next *Step) {
		for _, el := range array {
			inner.ExecNext(NewMessage(el), next)
		}
	}
}

// Source: trigger a Task execution for each element of the array.
func FromArray[T any](taskId string, array []T) *TTask[any, T] {
	return Via(Task[any](taskId), fromArray(array))
}

func fromString(string string, step ...int) Operator[any, string] {
	return func(inner *Inner, x *Message[any], next *Step) {
		subStr := ""

		for _, char := range strings.Split(string, "") {
			subStr += char

			if len(subStr) >= step[0] {
				inner.ExecNext(NewMessage(subStr), next)
				subStr = ""
			}
		}
	}
}

// Source: trigger a Task execution for each char of a string (or for each substring with a given step).
func FromString(taskId string, string string, step ...int) *TTask[any, string] {
	return Via(Task[any](taskId), fromString(string, step...))
}

func fromStringSplit(string string, delimiter string) Operator[any, string] {
	return func(inner *Inner, x *Message[any], next *Step) {
		for _, subStr := range strings.Split(string, delimiter) {
			inner.ExecNext(NewMessage(subStr), next)
		}
	}
}

// Source: trigger a Task execution for each substring, given a certain delimiter.
func FromStringSplit(taskId string, string string, delimiter string) *TTask[any, string] {
	return Via(Task[any](taskId), fromStringSplit(string, delimiter))
}

func fromInterval[T any](size time.Duration, max int, generator func(count int) T) Operator[any, T] {
	return func(inner *Inner, x *Message[any], next *Step) {
		counter := 0

		for range time.Tick(size) {
			value := generator(counter)

			inner.ExecNext(NewMessage(value), next)

			if max != 0 && counter == max-1 {
				break
			}

			counter++
		}
	}
}

// Source: trigger a task execution at a given interval.
// Generator function will produce the message, optionally using the interval counter.
// A max of 0 will generate and endless interval.
func FromInterval[T any](taskId string, size time.Duration, max int, generator func(count int) T) *TTask[any, T] {
	return Via(Task[any](taskId), fromInterval(size, max, generator))
}
