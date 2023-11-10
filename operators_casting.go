package ttask

func ToFloat64[T Number]() Operator[T, float64] {
	return func(inner *Inner, x *Message[T], step *Step) {
		inner.ExecNext(replaceValue(x, float64(x.Value)), step)
	}
}

func ToFloat32[T Number]() Operator[T, float32] {
	return func(inner *Inner, x *Message[T], step *Step) {
		inner.ExecNext(replaceValue(x, float32(x.Value)), step)
	}
}

func ToInt[T Number]() Operator[T, int] {
	return func(inner *Inner, x *Message[T], step *Step) {
		inner.ExecNext(replaceValue(x, int(x.Value)), step)
	}
}

func ToInt8[T Number]() Operator[T, int8] {
	return func(inner *Inner, x *Message[T], step *Step) {
		inner.ExecNext(replaceValue(x, int8(x.Value)), step)
	}
}

func ToInt16[T Number]() Operator[T, int16] {
	return func(inner *Inner, x *Message[T], step *Step) {
		inner.ExecNext(replaceValue(x, int16(x.Value)), step)
	}
}

func ToInt32[T Number]() Operator[T, int32] {
	return func(inner *Inner, x *Message[T], step *Step) {
		inner.ExecNext(replaceValue(x, int32(x.Value)), step)
	}
}

func ToInt64[T Number]() Operator[T, int32] {
	return func(inner *Inner, x *Message[T], step *Step) {
		inner.ExecNext(replaceValue(x, int64(x.Value)), step)
	}
}
