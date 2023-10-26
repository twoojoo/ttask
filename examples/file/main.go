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
	T(T(T(T(T(FromFile("./examples/file/in.txt"),
		Print[string]("#1 - received:\t\t")),
		Map(func(x string) string {
			return strings.ToUpper(x)
		})),
		Print[string]("#2 - transformed:\t")),
		ToFile("./examples/file/out.txt")),
		Print[string]("#3 - written:\t\t")).
		Run(context.Background())
}
