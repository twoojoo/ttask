package main

import (
	"context"
	"log"
	"strings"

	tt "github.com/twoojoo/ttask"
)

/**
 * NOTE: importing in an idiomatic way
 * makes everything too verbose
 * */

func main() {
	ctx := context.Background()

	in := "./examples/file/in.txt"
	out := "./examples/file/outt.txt"

	tt.T(tt.T(tt.T(tt.T(tt.T(
		tt.FromFile("t1", in),
		tt.Print[string]("#1 - received:\t\t")),
		tt.Map(func(x string) string {
			return strings.ToUpper(x)
		})),
		tt.Print[string]("#2 - transformed:\t")),
		tt.ToFile(out, "|")),
		tt.Print[string]("#3 - written:\t\t"),
	).Catch(func(i *tt.Inner, e error) {
		log.Fatal(e)
	}).Run(ctx)
}
