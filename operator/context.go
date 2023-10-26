package operator

import (
	"context"

	"github.com/twoojoo/ttask/task"
)

//Cache a key/value record in the Task context. Use an extractor function to pull the value from the processed item.
func WithContextValue[T any](k any, ext func(x T) any) task.Operator[T, T] {
	return func(m *task.Meta, x *task.Message[T], next *task.Step) {
		m.Ctx = context.WithValue(m.Ctx, k, ext(x.Value))
		task.ExecNext(m, x, next)
	}
}
