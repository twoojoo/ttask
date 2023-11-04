package ttask

import (
	"io"
)

func fromReader(r io.Reader, bufSize int) Operator[any, []byte] {
	return func(inner *Inner, x *Message[any], next *Step) {
		buf := make([]byte, bufSize)

		for {
			n, err := r.Read(buf)

			if err != io.EOF {
				if err != nil {
					inner.Error(err)
					return
				}

				inner.ExecNext(newMessage(buf[:n]), next)
			}
		}
	}
}

func FromReader(taskId string, r io.Reader, bufSize int) *TTask[any, []byte] {
	return Via(Task[any](taskId), fromReader(r, bufSize))
}
