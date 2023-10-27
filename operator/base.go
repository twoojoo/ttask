package operator

import (
	"fmt"

	"github.com/twoojoo/ttask/task"
)

func Print[T any](prefix ...string) task.Operator[T, T] {
	return func(m *task.Meta, x *task.Message[T], step *task.Step) {
		if len(prefix) > 0 {
			fmt.Println(prefix[0], x.Value)
		} else {
			fmt.Println(x.Value)
		}

		m.ExecNext(x, step)
	}
}

func PrintRaw[T any](prefix ...string) task.Operator[T, T] {
	return func(m *task.Meta, x *task.Message[T], step *task.Step) {
		if len(prefix) > 0 {
			fmt.Printf("%s %+v\n", prefix[0], x)
		} else {
			fmt.Printf("%+v\n", x)
		}

		m.ExecNext(x, step)
	}
}

func Map[T, R any](cb func(x T) R) task.Operator[T, R] {
	return func(m *task.Meta, x *task.Message[T], next *task.Step) {
		m.ExecNext(task.ReplaceValue(x, cb(x.Value)), next)
	}
}

func MapRaw[T, R any](cb func(m *task.Meta, x *task.Message[T]) R) task.Operator[T, R] {
	return func(m *task.Meta, x *task.Message[T], next *task.Step) {
		m.ExecNext(task.ReplaceValue(x, task.ReplaceValue(x, cb(m, x))), next)
	}
}

func Filter[T, R any](cb func(x T) bool) task.Operator[T, T] {
	return func(m *task.Meta, x *task.Message[T], next *task.Step) {
		ok := cb(x.Value)
		if ok {
			m.ExecNext(x, next)
		}
	}
}

func FilterRaw[T, R any](cb func(m *task.Meta, x *task.Message[T]) bool) task.Operator[T, T] {
	return func(m *task.Meta, x *task.Message[T], next *task.Step) {
		ok := cb(m, x)
		if ok {
			m.ExecNext(x, next)
		}
	}
}

func Tap[T any](cb func(x T)) task.Operator[T, T] {
	return func(m *task.Meta, x *task.Message[T], next *task.Step) {
		cb(x.Value)
		m.ExecNext(x, next)
	}
}

func TapRaw[T any](cb func(m *task.Meta, x *task.Message[T])) task.Operator[T, T] {
	return func(m *task.Meta, x *task.Message[T], next *task.Step) {
		cb(m, x)
		m.ExecNext(x, next)
	}
}
