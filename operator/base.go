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

		task.ExecNext(m, x, step)
	}
}

func Map[T, R any](cb func(x T) R) task.Operator[T, R] {
	return func(m *task.Meta, x *task.Message[T], next *task.Step) {
		task.ExecNext(m, task.ReplaceValue(x, cb(x.Value)), next)
	}
}

func MapRaw[T, R any](cb func(m *task.Meta, x *task.Message[T]) R) task.Operator[T, R] {
	return func(m *task.Meta, x *task.Message[T], next *task.Step) {
		task.ExecNext(m, task.ReplaceValue(x, task.ReplaceValue(x, cb(m, x))), next)
	}
}

func Filter[T, R any](cb func(x T) bool) task.Operator[T, T] {
	return func(m *task.Meta, x *task.Message[T], next *task.Step) {
		ok := cb(x.Value)
		if ok {
			task.ExecNext(m, x, next)
		}
	}
}

func FilterRaw[T, R any](cb func(m *task.Meta, x *task.Message[T]) bool) task.Operator[T, T] {
	return func(m *task.Meta, x *task.Message[T], next *task.Step) {
		ok := cb(m, x)
		if ok {
			task.ExecNext(m, x, next)
		}
	}
}

func Tap[T any](cb func(x T)) task.Operator[T, T] {
	return func(m *task.Meta, x *task.Message[T], next *task.Step) {
		cb(x.Value)
		task.ExecNext(m, x, next)
	}
}

func TapRaw[T any](cb func(m *task.Meta, x *task.Message[T])) task.Operator[T, T] {
	return func(m *task.Meta, x *task.Message[T], next *task.Step) {
		cb(m, x)
		task.ExecNext(m, x, next)
	}
}
