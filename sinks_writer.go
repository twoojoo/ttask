package ttask

import (
	"io"
)

// Sink: write message to a writer
func ToWriter[T any](w io.Writer, toBytes func(x T) []byte) Operator[T, T] {
	return MapRaw(func(inner *Inner, x *Message[T]) T {
		_, err := w.Write(toBytes(x.Value))

		if err != nil {
			inner.Error(err)
			return x.Value
		}

		return x.Value
	})
}

// Sink: write message to a writer. Next message value will be the number of written bytes.
func ToWriterCount[T any](w io.Writer, toBytes func(x T) []byte) Operator[T, int] {
	return MapRaw(func(inner *Inner, x *Message[T]) int {
		w, err := w.Write(toBytes(x.Value))

		if err != nil {
			inner.Error(err)
			return 0
		}

		return w
	})
}
