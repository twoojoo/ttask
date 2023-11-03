package sink

import (
	"io"

	"github.com/twoojoo/ttask/operator"
	"github.com/twoojoo/ttask/task"
)

// Sink: write message to a writer
func ToWriter[T any](w io.Writer, toBytes func(x T) []byte) task.Operator[T, T] {
	return operator.MapRaw(func(m *task.Inner, x *task.Message[T]) T {
		_, err := w.Write(toBytes(x.Value))

		if err != nil {
			m.Error(err)
			return x.Value
		}
		
		return x.Value
	})
}


// Sink: write message to a writer. Next message value will be the number of written bytes.
func ToWriterCount[T any](w io.Writer, toBytes func(x T) []byte) task.Operator[T, int] {
	return operator.MapRaw(func(m *task.Inner, x *task.Message[T]) int {
		w, err := w.Write(toBytes(x.Value))

		if err != nil {
			m.Error(err)
			return 0
		}

		return w
	})
}
