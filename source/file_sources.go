package operator

import (
	"bufio"
	"os"

	"github.com/twoojoo/ttask/task"
)

func fromFile(path string) task.Operator[any, string] {
	return func(inner *task.Inner, _ *task.Message[any], next *task.Step) {
		file, err := os.Open(path)

		if err != nil {
			inner.Error(err)
			return
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()
			inner.ExecNext(task.NewMessage(line), next)
		}

		if err := scanner.Err(); err != nil {
			inner.Error(err)
			return
		}
	}
}

// Source: read a file an trigger a Task execution for each line.
func FromFile(taskId string, path string) *task.TTask[any, string] {
	return task.T(task.Task[any](taskId), fromFile(path))
}
