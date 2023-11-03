package main

import (
	"context"
	"log"
	"strings"

	. "github.com/twoojoo/ttask/operator"
	. "github.com/twoojoo/ttask/sink"
	. "github.com/twoojoo/ttask/source"
	. "github.com/twoojoo/ttask/task"
)

func main() {
	ctx := context.Background()

	in := "./examples/file/in.txt"
	out := "./examples/file/out.txt"

	T(T(T(T(T(
		FromFile("t1", in),
		Print[string]("#1 - received:\t\t")),
		Map(func(x string) string {
			return strings.ToUpper(x)
		})),
		Print[string]("#2 - transformed:\t")),
		ToFile(out, "|")),
		Print[string]("#3 - written:\t\t"),
	).Catch(func(m *Inner, e error) {
		log.Fatal(e)
	}).Run(ctx)
}
