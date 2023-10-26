package main

import (
	"context"
	"strings"

	. "github.com/twoojoo/ttask/operator"
	. "github.com/twoojoo/ttask/source"
	. "github.com/twoojoo/ttask/task"
	. "github.com/twoojoo/ttask/sink"
)

func main() {
	ctx := context.Background()

	in := "./examples/file/out.txt"
	out := "./examples/file/in.txt"

	T(T(T(T(T(FromFile(in),
		Print[string]("#1 - received:\t\t")),
		Map(func(x string) string {
			return strings.ToUpper(x)
		})),
		Print[string]("#2 - transformed:\t")),
		ToFile(out)),
		Print[string]("#3 - written:\t\t")).
		Run(ctx)
}
