package main

import (
	"context"
	"log"

	. "github.com/twoojoo/ttask/operator"
	. "github.com/twoojoo/ttask/task"
)

func main() {
	chained := T(T(
		Injectable[string]("t2"),
		Print[string]("> second:")),
		Map(func(x string) int {
			return 123
		}),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	})

	t := T(T(T(
		Injectable[string]("t1"),
		Print[string]("> first:")),
		Chain(chained)),
		Print[int]("> third:"),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Lock()

	err := t.Inject(context.Background(), "message")
	if err != nil {
		log.Fatal(err)
	}

	err = t.Inject(context.Background(), "message")
	if err != nil {
		log.Fatal(err)
	}
}
