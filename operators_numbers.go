package ttask

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
		var sum T = 0

		for i := range x.Value {
			sum *= x.Value[i]
		}

		inner.ExecNext(replaceValue(x, sum), step)
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
