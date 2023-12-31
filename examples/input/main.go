package main

import (
	"context"
	"log"

	. "github.com/twoojoo/ttask"
)

func main() {
	ctx := context.Background()

	T(
		FromReadline("t1", "> type a phrase: "),
		Print[string]("> you typed:"),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Run(ctx)

	T(T(
		FromReadChar("t2", "> now type a char: "),
		Map(func(x rune) string { return string(x) })),
		Print[string]("> you typed:"),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Run(ctx)

}
