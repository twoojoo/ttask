package operator

import (
	"context"

	"github.com/twoojoo/ttask/task"
)

// Cache a key/value record in the Task context. Use an extractor function to pull the value from the processed item.
func WithContextValue[T any](k any, ext func(x T) any) task.Operator[T, T] {
	return func(inner *task.Inner, x *task.Message[T], next *task.Step) {
		inner.Context = context.WithValue(inner.Context, k, ext(x.Value))
		inner.ExecNext(x, next)
	}
}
