package ttask

import (
	"bufio"
	"io"
	"os"
)

func FromReader(taskId string, r io.Reader, bufSize int) *TTask[any, []byte] {
	return Via(Task[any](taskId), fromReader(r, bufSize))
}

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

				inner.ExecNext(NewMessage(buf[:n]), next)
			}
		}
	}
}

func FromStdin(taskId string) *TTask[any, string] {
	return Via(Task[any](taskId), fromStdin())
}

func fromStdin() Operator[any, string] {
	return func(inner *Inner, x *Message[any], next *Step) {
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			inner.ExecNext(NewMessage(scanner.Text()), next)
		}
	}
}
