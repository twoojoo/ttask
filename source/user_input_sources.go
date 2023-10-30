package operator

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/twoojoo/ttask/task"
)

func fromReadLine(prompt string) task.Operator[any, string] {
	return func(m *task.Meta, x *task.Message[any], next *task.Step) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print(prompt)
		line, err := reader.ReadString('\n')

		if err != nil {
			m.Error(err)
			return
		}

		line = strings.Split(line, "\n")[0]

		m.ExecNext(task.NewMessage(line), next)
	}
}

func FromReadline(taskId string, prompt string) *task.TTask[any, string] {
	return task.T(task.Task[any](taskId), fromReadLine(prompt))
}

func fromReadChar(prompt string) task.Operator[any, rune] {
	return func(m *task.Meta, x *task.Message[any], next *task.Step) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print(prompt)
		rune, _, err := reader.ReadRune()
		if err != nil {
			m.Error(err)
			return
		}

		m.ExecNext(task.NewMessage(rune), next)
	}
}

func FromReadChar(taskId string, prompt string) *task.TTask[any, rune] {
	return task.T(task.Task[any](taskId), fromReadChar(prompt))
}
