package ttask

import (
	"io"
)

// Sink: write each Task result to a file unsing a separator (default: \n)
func ToFile(path string, separator ...string) Operator[string, string] {
	return func(inner *Inner, x *Message[string], next *Step) {
		file, err := openOrCreateFile(path)
		defer file.Close()

		if err != nil {
			inner.Error(err)
			return
		}

		s := "\n"
		if len(separator) > 0 {
			s = separator[0]
		}

		writer := io.StringWriter(file)

		_, err = writer.WriteString(x.Value + s)
		if err != nil {
			inner.Error(err)
			return
		}

		inner.ExecNext(x, next)
	}
}
