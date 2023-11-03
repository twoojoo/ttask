package sink

import (
	"io"

	"github.com/twoojoo/ttask/task"
	"github.com/twoojoo/ttask/utils"
)

//Sink: write each Task result to a file unsing a separator (default: \n)
func ToFile(path string, separator ...string) task.Operator[string, string] {
	return func(m *task.Inner, x *task.Message[string], next *task.Step) {
		file, err := utils.OpenOrCreateFile(path)
		defer file.Close()

		if err != nil {
			m.Error(err)
			return
		}

		s := "\n"
		if len(separator) > 0 {
			s = separator[0]
		}

		writer := io.StringWriter(file)

		_, err = writer.WriteString(x.Value + s)
		if err != nil {
			m.Error(err)
			return
		}


		m.ExecNext(x, next)
	}
}
