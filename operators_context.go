package ttask

import (
	"context"
)

// Cache a key/value record in the Task context. Use an extractor function to pull the value from the processed item.
func WithContextValue[T any](k any, ext func(x T) any) Operator[T, T] {
	return func(inner *Inner, x *Message[T], next *Step) {
		inner.Context = context.WithValue(inner.Context, k, ext(x.Value))
		inner.ExecNext(x, next)
	}
}
