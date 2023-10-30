package operator

import (
	"strings"
	"time"

	"github.com/twoojoo/ttask/task"
)

func fromItem[T any](item T) task.Operator[any, T] {
	return func(m *task.Meta, x *task.Message[any], next *task.Step) {
		m.ExecNext(task.NewMessage(item), next)
	}
}

// Source: trigger a task once with the give item.
func FromItem[T any](taskId string, item T) *task.TTask[any, T] {
	return task.T(task.Task[any](taskId), fromItem(item))
}

func fromArray[T any](array []T) task.Operator[any, T] {
	return func(m *task.Meta, x *task.Message[any], next *task.Step) {
		for _, el := range array {
			m.ExecNext(task.NewMessage(el), next)
		}
	}
}

// Source: trigger a Task execution for each element of the array.
func FromArray[T any](taskId string, array []T) *task.TTask[any, T] {
	return task.T(task.Task[any](taskId), fromArray(array))
}

func fromString(string string, step ...int) task.Operator[any, string] {
	return func(m *task.Meta, x *task.Message[any], next *task.Step) {
		subStr := ""

		for _, char := range strings.Split(string, "") {
			subStr += char

			if len(subStr) >= step[0] {
				m.ExecNext(task.NewMessage(subStr), next)
				subStr = ""
			}
		}
	}
}

// Source: trigger a Task execution for each char of a string (or for each substring with a given step).
func FromString(taskId string, string string, step ...int) *task.TTask[any, string] {
	return task.T(task.Task[any](taskId), fromString(string, step...))
}

func fromStringSplit(string string, delimiter string) task.Operator[any, string] {
	return func(m *task.Meta, x *task.Message[any], next *task.Step) {
		for _, subStr := range strings.Split(string, delimiter) {
			m.ExecNext(task.NewMessage(subStr), next)
		}
	}
}

// Source: trigger a Task execution for each substring, given a certain delimiter.
func FromStringSplit(taskId string, string string, delimiter string) *task.TTask[any, string] {
	return task.T(task.Task[any](taskId), fromStringSplit(string, delimiter))
}

func fromInterval[T any](size time.Duration, max int, generator func(count int) T)  task.Operator[any, T] {
	return func(m *task.Meta, x *task.Message[any], next *task.Step) {
		counter := 0

		for range time.Tick(size) {
			value  := generator(counter)

			m.ExecNext(task.NewMessage(value), next)

			if counter == max-1 {
				break
			}

			counter++
		}
	}
}

// Source: trigger a task execution at a given interval. 
// Generator function will produce the message, optionally using the interval counter.
func FromInterval[T any](taskId string, size time.Duration, max int, generator func(count int) T) *task.TTask[any, T] {
	return task.T(task.Task[any](taskId), fromInterval(size, max, generator))
}