package main

import (
	"context"
	"log"
	"strings"

	"github.com/twoojoo/ttask"
)

func main() {
	ctx := context.Background()

	in := "./examples/file/in.txt"
	out := "./examples/file/out.txt"

	ttask.T(ttask.T(ttask.T(ttask.T(ttask.T(
		ttask.FromFile("t1", in),
		ttask.Print[string]("#1 - received:\t\t")),
		ttask.Map(func(x string) string {
			return strings.ToUpper(x)
		})),
		ttask.Print[string]("#2 - transformed:\t")),
		ttask.ToFile(out, "|")),
		ttask.Print[string]("#3 - written:\t\t"),
	).Catch(func(i *ttask.Inner, e error) {
		log.Fatal(e)
	}).Run(ctx)
}
