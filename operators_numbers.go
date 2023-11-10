package ttask

import "golang.org/x/exp/constraints"

func Sum[T Number]() Operator[[]T, T] {
	return func(inner *Inner, x *Message[[]T], step *Step) {
		var sum T = 0

		for i := range x.Value {
			sum += x.Value[i]
		}

		inner.ExecNext(replaceValue(x, sum), step)
	}
}

func Multiply[T Number]() Operator[[]T, T] {
	return func(inner *Inner, x *Message[[]T], step *Step) {
		var product T = 1

		for i := range x.Value {
			product = product * x.Value[i]
		}

		inner.ExecNext(replaceValue(x, product), step)
	}
}

func Min[T Number]() Operator[[]T, T] {
	return func(inner *Inner, x *Message[[]T], step *Step) {
		var min T

		for i := range x.Value {
			if i == 0 {
				min = x.Value[i]
			} else if min > x.Value[i] {
				min = x.Value[i]
			}
		}

		inner.ExecNext(replaceValue(x, min), step)
	}
}

func Max[T Number]() Operator[[]T, T] {
	return func(inner *Inner, x *Message[[]T], step *Step) {
		var min T

		for i := range x.Value {
			if i == 0 {
				min = x.Value[i]
			} else if min > x.Value[i] {
				min = x.Value[i]
			}
		}

		inner.ExecNext(replaceValue(x, min), step)
	}
}

func Average[X Number]() Operator[[]X, X] {
	return func(inner *Inner, x *Message[[]X], step *Step) {
		var sum X = 0

		for i := range x.Value {
			sum += x.Value[i]
		}

		l := X(len(x.Value))

		inner.ExecNext(replaceValue(x, sum/l), step)
	}
}

func IsEven[T constraints.Integer]() func(x T) bool {
	return func(x T) bool {
		return x%2 == 0
	}
}

func IsOdd[T constraints.Integer](x int) bool {
	return x%2 == 1
}
