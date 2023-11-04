package ttask

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func fromReadLine(prompt string) Operator[any, string] {
	return func(inner *Inner, x *Message[any], next *Step) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print(prompt)
		line, err := reader.ReadString('\n')

		if err != nil {
			inner.Error(err)
			return
		}

		line = strings.Split(line, "\n")[0]

		inner.ExecNext(NewMessage(line), next)
	}
}

func FromReadline(taskId string, prompt string) *TTask[any, string] {
	return Via(Task[any](taskId), fromReadLine(prompt))
}

func fromReadChar(prompt string) Operator[any, rune] {
	return func(inner *Inner, x *Message[any], next *Step) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print(prompt)
		rune, _, err := reader.ReadRune()
		if err != nil {
			inner.Error(err)
			return
		}

		inner.ExecNext(NewMessage(rune), next)
	}
}

func FromReadChar(taskId string, prompt string) *TTask[any, rune] {
	return Via(Task[any](taskId), fromReadChar(prompt))
}
