package operator

import (
	"github.com/twoojoo/ttask/task"
)

func MapArray[T, R any](cb func(x T) R) task.Operator[[]T, []R] {
	return func(inner *task.Inner, x *task.Message[[]T], next *task.Step) {
		mapped := make([]R, len(x.Value))

		for i := 0; i < len(x.Value); i++ {
			mapped[i] = cb(x.Value[i])
		}

		inner.ExecNext(task.ReplaceValue(x, mapped), next)
	}
}

func MapArrayRaw[T, R any](cb func(inner *task.Inner, x *task.Message[T]) R) task.Operator[[]T, []R] {
	return func(inner *task.Inner, x *task.Message[[]T], next *task.Step) {
		mapped := make([]R, len(x.Value))

		for i := 0; i < len(x.Value); i++ {
			mapped[i] = cb(inner, task.ReplaceValue(x, x.Value[i]))
		}

		inner.ExecNext(task.ReplaceValue(x, mapped), next)
	}
}

func ForEach[T any](cb func(x T)) task.Operator[[]T, []T] {
	return func(inner *task.Inner, x *task.Message[[]T], next *task.Step) {
		for i := 0; i < len(x.Value); i++ {
			cb(x.Value[i])
		}

		inner.ExecNext(x, next)
	}
}

func EachRaw[T, R any](cb func(inner *task.Inner, x *task.Message[T]) R) task.Operator[[]T, []R] {
	return func(inner *task.Inner, x *task.Message[[]T], next *task.Step) {
		for i := 0; i < len(x.Value); i++ {
			cb(inner, task.ReplaceValue(x, x.Value[i]))
		}

		inner.ExecNext(x, next)
	}
}

func FilterArray[T any](cb func(x T) bool) task.Operator[[]T, []T] {
	return func(inner *task.Inner, x *task.Message[[]T], next *task.Step) {
		filtered := []T{}

		for i := 0; i < len(x.Value); i++ {
			if cb(x.Value[i]) {
				filtered = append(filtered, x.Value[i])
			}
		}

		inner.ExecNext(task.ReplaceValue(x, filtered), next)
	}
}

func FilterArrayRaw[T any](cb func(inner *task.Inner, x *task.Message[T]) bool) task.Operator[[]T, []T] {
	return func(inner *task.Inner, x *task.Message[[]T], next *task.Step) {
		filtered := []T{}

		for i := 0; i < len(x.Value); i++ {
			if cb(inner, task.ReplaceValue(x, x.Value[i])) {
				filtered = append(filtered, x.Value[i])
			}
		}

		inner.ExecNext(task.ReplaceValue(x, filtered), next)
	}
}

// JS-like array reducer
func ReduceArray[T, R any](base R, reducer func(acc *R, x T) R) task.Operator[[]T, R] {
	return func(inner *task.Inner, x *task.Message[[]T], next *task.Step) {
		for i := 0; i < len(x.Value); i++ {
			base = reducer(&base, x.Value[i])
		}

		inner.ExecNext(task.ReplaceValue(x, base), next)
	}
}

// JS-like array reducer [raw version]
func ReduceArrayRaw[T, R any](base R, reducer func(acc *R, inner *task.Inner, x *task.Message[T]) R) task.Operator[[]T, R] {
	return func(inner *task.Inner, x *task.Message[[]T], next *task.Step) {
		for i := 0; i < len(x.Value); i++ {
			base = reducer(&base, inner, task.ReplaceValue(x, x.Value[i]))
		}

		inner.ExecNext(task.ReplaceValue(x, base), next)
	}
}

func FlatArray[T any]() task.Operator[[][]T, []T] {
	return func(inner *task.Inner, x *task.Message[[][]T], next *task.Step) {
		flatten := flatArray(&x.Value)
		inner.ExecNext(task.ReplaceValue(x, flatten), next)
	}
}

func flatArray[T any](arr *[][]T) []T {
	flatten := []T{}

	for i := range *arr {
		flatten = append(flatten, (*arr)[i]...)
	}

	return flatten
}

// Continue the task execution for each element of the array synchronously
func IterateArray[T any]() task.Operator[[]T, T] {
	return func(inner *task.Inner, x *task.Message[[]T], next *task.Step) {
		for i := 0; i < len(x.Value); i++ {
			inner.ExecNext(task.ReplaceValue(x, x.Value[i]), next)
		}
	}
}

// Continue the task exection for each element of the array asynchronously
func ParallelizeArray[T any]() task.Operator[[]T, T] {
	return func(inner *task.Inner, x *task.Message[[]T], next *task.Step) {
		ch := make(chan struct{}, len(x.Value))

		for i := 0; i < len(x.Value); i++ {
			c := *&i
			go func() {
				inner.ExecNext(task.ReplaceValue(x, x.Value[c]), next)
				ch <- struct{}{}
			}()
		}

		//wait for all the iterations to complete
		for i := 0; i < len(x.Value); i++ {
			<-ch
		}
	}
}
