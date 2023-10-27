package operator

import (
	"strings"

	"github.com/twoojoo/ttask/task"
)

func fromItem[T any](item T) task.Operator[any, T] {
	return func(m *task.Meta, x *task.Message[any], next *task.Step) {
		m.ExecNext(task.NewMessage(item), next)
	}
}

//Source: trigger a task once with the give item.
func FromItem[T any](item T) *task.TTask[any, T] {
	return task.T(task.Task[any](), fromItem(item))
}

func fromArray[T any](array []T) task.Operator[any, T] {
	return func(m *task.Meta, x *task.Message[any], next *task.Step) {
		for _, el := range array {
			m.ExecNext(task.NewMessage(el), next)
		}
	}
}

//Source: trigger a Task execution for each element of the array.
func FromArray[T any](array []T) *task.TTask[any, T] {
	return task.T(task.Task[any](), fromArray(array))
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

//Source: trigger a Task execution for each char of a string (or for each substring with a given step).
func FromString(string string, step ...int) *task.TTask[any, string] {
	return task.T(task.Task[any](), fromString(string, step...))
}

func fromStringSplit(string string, delimiter string) task.Operator[any, string] {
	return func(m *task.Meta, x *task.Message[any], next *task.Step) {
		for _, subStr := range strings.Split(string, delimiter) {
			m.ExecNext(task.NewMessage(subStr), next)
		}
	}
}

//Source: trigger a Task execution for each substring, given a certain delimiter.
func FromStringSplit(string string, delimiter string) *task.TTask[any, string] {
	return task.T(task.Task[any](), fromStringSplit(string, delimiter))
}
