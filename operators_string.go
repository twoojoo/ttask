package ttask

func Length() Operator[string, int] {
	return func(inner *Inner, x *Message[string], step *Step) {
		inner.ExecNext(replaceValue(x, len(x.Value)), step)
	}
}

func Concat(separator ...string) Operator[[]string, string] {
	return func(inner *Inner, x *Message[[]string], step *Step) {
		sep := ""
		if len(separator) > 0 {
			sep = separator[0]
		}

		result := ""
		for i := range x.Value {
			result += x.Value[i]
			if i != len(x.Value)-1 {
				result += sep
			}
		}

		inner.ExecNext(replaceValue(x, result), step)
	}
}
