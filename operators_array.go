package ttask

func MapArray[T, R any](cb func(x T) R) Operator[[]T, []R] {
	return func(inner *Inner, x *Message[[]T], next *Step) {
		mapped := make([]R, len(x.Value))

		for i := 0; i < len(x.Value); i++ {
			mapped[i] = cb(x.Value[i])
		}

		inner.ExecNext(replaceValue(x, mapped), next)
	}
}

func MapArrayRaw[T, R any](cb func(inner *Inner, x *Message[T]) R) Operator[[]T, []R] {
	return func(inner *Inner, x *Message[[]T], next *Step) {
		mapped := make([]R, len(x.Value))

		for i := 0; i < len(x.Value); i++ {
			mapped[i] = cb(inner, replaceValue(x, x.Value[i]))
		}

		inner.ExecNext(replaceValue(x, mapped), next)
	}
}

func ForEach[T any](cb func(x T)) Operator[[]T, []T] {
	return func(inner *Inner, x *Message[[]T], next *Step) {
		for i := 0; i < len(x.Value); i++ {
			cb(x.Value[i])
		}

		inner.ExecNext(x, next)
	}
}

func EachRaw[T, R any](cb func(inner *Inner, x *Message[T]) R) Operator[[]T, []R] {
	return func(inner *Inner, x *Message[[]T], next *Step) {
		for i := 0; i < len(x.Value); i++ {
			cb(inner, replaceValue(x, x.Value[i]))
		}

		inner.ExecNext(x, next)
	}
}

func FilterArray[T any](cb func(x T) bool) Operator[[]T, []T] {
	return func(inner *Inner, x *Message[[]T], next *Step) {
		filtered := []T{}

		for i := 0; i < len(x.Value); i++ {
			if cb(x.Value[i]) {
				filtered = append(filtered, x.Value[i])
			}
		}

		inner.ExecNext(replaceValue(x, filtered), next)
	}
}

func FilterArrayRaw[T any](cb func(inner *Inner, x *Message[T]) bool) Operator[[]T, []T] {
	return func(inner *Inner, x *Message[[]T], next *Step) {
		filtered := []T{}

		for i := 0; i < len(x.Value); i++ {
			if cb(inner, replaceValue(x, x.Value[i])) {
				filtered = append(filtered, x.Value[i])
			}
		}

		inner.ExecNext(replaceValue(x, filtered), next)
	}
}

// JS-like array reducer
func ReduceArray[T, R any](base R, reducer func(acc *R, x T) R) Operator[[]T, R] {
	baseCp := *&base

	return func(inner *Inner, x *Message[[]T], next *Step) {
		for i := 0; i < len(x.Value); i++ {
			baseCp = reducer(&baseCp, x.Value[i])
		}

		inner.ExecNext(replaceValue(x, baseCp), next)

		baseCp = base
	}
}

// JS-like array reducer [raw version]
func ReduceArrayRaw[T, R any](base R, reducer func(acc *R, inner *Inner, x *Message[T]) R) Operator[[]T, R] {
	baseCp := *&base

	return func(inner *Inner, x *Message[[]T], next *Step) {

		for i := 0; i < len(x.Value); i++ {
			baseCp = reducer(&baseCp, inner, replaceValue(x, x.Value[i]))
		}

		inner.ExecNext(replaceValue(x, baseCp), next)

		baseCp = base
	}
}

func FlatArray[T any]() Operator[[][]T, []T] {
	return func(inner *Inner, x *Message[[][]T], next *Step) {
		flatten := flatArray(&x.Value)
		inner.ExecNext(replaceValue(x, flatten), next)
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
func IterateArray[T any]() Operator[[]T, T] {
	return func(inner *Inner, x *Message[[]T], next *Step) {
		for i := 0; i < len(x.Value); i++ {
			inner.ExecNext(replaceValue(x, x.Value[i]), next)
		}
	}
}

// Continue the task exection for each element of the array asynchronously
func ParallelizeArray[T any]() Operator[[]T, T] {
	return func(inner *Inner, x *Message[[]T], next *Step) {
		ch := make(chan struct{}, len(x.Value))

		for i := 0; i < len(x.Value); i++ {
			c := *&i
			go func() {
				inner.ExecNext(replaceValue(x, x.Value[c]), next)
				ch <- struct{}{}
			}()
		}

		//wait for all the iterations to complete
		for i := 0; i < len(x.Value); i++ {
			<-ch
		}
	}
}
