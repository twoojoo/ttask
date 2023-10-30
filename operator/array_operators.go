package operator

import (
	"github.com/twoojoo/ttask/task"
)

func MapArray[T, R any](cb func(x T) R) task.Operator[[]T, []R] {
	return func(m *task.Meta, x *task.Message[[]T], next *task.Step) {
		mapped := make([]R, len(x.Value))

		for i := 0; i < len(x.Value); i++ {
			mapped[i] = cb(x.Value[i])
		}

		m.ExecNext(task.ReplaceValue(x, mapped), next)
	}
}

func MapArrayRaw[T, R any](cb func(m *task.Meta, x T) R) task.Operator[[]T, []R] {
	return func(m *task.Meta, x *task.Message[[]T], next *task.Step) {
		mapped := make([]R, len(x.Value))

		for i := 0; i < len(x.Value); i++ {
			mapped[i] = cb(m, x.Value[i])
		}

		m.ExecNext(task.ReplaceValue(x, mapped), next)
	}
}

func Each[T any](cb func(x T)) task.Operator[[]T, []T] {
	return func(m *task.Meta, x *task.Message[[]T], next *task.Step) {
		for i := 0; i < len(x.Value); i++ {
			cb(x.Value[i])
		}

		m.ExecNext(x, next)
	}
}

func EachRaw[T, R any](cb func(m *task.Meta, x T) R) task.Operator[[]T, []R] {
	return func(m *task.Meta, x *task.Message[[]T], next *task.Step) {
		for i := 0; i < len(x.Value); i++ {
			cb(m, x.Value[i])
		}

		m.ExecNext(x, next)
	}
}

func FilterArray[T any](cb func(x T) bool) task.Operator[[]T, []T] {
	return func(m *task.Meta, x *task.Message[[]T], next *task.Step) {
		filtered := []T{}

		for i := 0; i < len(x.Value); i++ {
			if cb(x.Value[i]) {
				filtered = append(filtered, x.Value[i])
			}
		}

		m.ExecNext(task.ReplaceValue(x, filtered), next)
	}
}

func FilterArrayRaw[T any](cb func(m *task.Meta, x T) bool) task.Operator[[]T, []T] {
	return func(m *task.Meta, x *task.Message[[]T], next *task.Step) {
		filtered := []T{}

		for i := 0; i < len(x.Value); i++ {
			if cb(m, x.Value[i]) {
				filtered = append(filtered, x.Value[i])
			}
		}

		m.ExecNext(task.ReplaceValue(x, filtered), next)
	}
}

// func FlatArray[T any](cb func(x T) bool) task.Operator[[]T, []T] {
// 	return func(m *task.Meta, x *task.Message[[]T], next *task.Step) {
// 		filtered := []T{}

// 		for i := 0; i < len(x.Value); i++ {
// 			if cb(x.Value[i]) {
// 				filtered = append(filtered, x.Value[i])
// 			}
// 		}

// 		m.ExecNext(task.ReplaceValue(x, filtered), next)
// 	}
// }

// func FlatArrayRaw[T any](cb func(m *task.Meta, x T) bool) task.Operator[[][]T, []T] {
// 	return func(m *task.Meta, x *task.Message[[][]T], next *task.Step) {
// 		filtered := []T{}

// 		for i := 0; i < len(x.Value); i++ {
// 			if cb(m, x.Value[i]) {
// 				filtered = append(filtered, x.Value[i])
// 			}
// 		}

// 		m.ExecNext(task.ReplaceValue(x, filtered), next)
// 	}
// }