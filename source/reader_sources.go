package operator

import (
	"io"

	"github.com/twoojoo/ttask/task"
)

func fromReader(r io.Reader, bufSize int) task.Operator[any, []byte] {
	return func(m *task.Inner, x *task.Message[any], next *task.Step) {
		buf := make([]byte, bufSize)

		for {
			n, err := r.Read(buf)

			if err != io.EOF {
				if err != nil {
					m.Error(err)
					return
				}

				m.ExecNext(task.NewMessage(buf[:n]), next)
			}
		}
	}
}

func FromReader(taskId string, r io.Reader, bufSize int) *task.TTask[any, []byte] {
	return task.T(task.Task[any](taskId), fromReader(r, bufSize))
}
