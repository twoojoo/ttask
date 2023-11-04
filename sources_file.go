package ttask

import (
	"bufio"
	"os"
)

func fromFile(path string) Operator[any, string] {
	return func(inner *Inner, _ *Message[any], next *Step) {
		file, err := os.Open(path)

		if err != nil {
			inner.Error(err)
			return
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()
			inner.ExecNext(newMessage(line), next)
		}

		if err := scanner.Err(); err != nil {
			inner.Error(err)
			return
		}
	}
}

// Source: read a file an trigger a Task execution for each line.
func FromFile(taskId string, path string) *TTask[any, string] {
	return Via(Task[any](taskId), fromFile(path))
}
