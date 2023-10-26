package operator

import (
	"bufio"
	"os"

	"github.com/twoojoo/ttask/task"
)

func fromFile(path string) task.Operator[any, string] {
	return func(m *task.Meta, _ *task.Message[any], next *task.Step) {
		file, err := os.Open(path)

		if err != nil {
			panic(err)
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()
			task.ExecNext(m, task.NewMessage(line), next)
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}
}

func FromFile(path string) *task.TTask[any, string] {
	return task.T(task.Task[any](), fromFile(path))
}
