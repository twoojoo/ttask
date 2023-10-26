package sink

import (
	"io"

	"github.com/twoojoo/ttask/task"
	"github.com/twoojoo/ttask/utils"
)

func ToFile(path string) task.Operator[string, string] {
	return func(m *task.Meta, x *task.Message[string], next *task.Step) {
		file, err := utils.OpenOrCreateFile(path)
		if err != nil {
			panic(err)
		}

		writer := io.StringWriter(file)

		_, err = writer.WriteString(x.Value + "\n")
		if err != nil {
			panic(err)
		}

		defer file.Close()

		task.ExecNext(m, x, next)
	}
}
