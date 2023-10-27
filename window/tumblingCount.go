package window

import (
	"github.com/twoojoo/ttask/storage"
	"github.com/twoojoo/ttask/task"
)

func TumblingWindowCount[T any](storage storage.Storage[task.Message[T]], maxSize int) task.Operator[T, []T] {
	storage.Init()

	return func(m *task.Meta, x *task.Message[T], next *task.Step) {
		size := storage.Push(x.Key, x)

		if size >= maxSize {
			items := storage.Flush(x.Key)
			m.ExecNext(task.ToArray(x, items), next)
		}
	}
}