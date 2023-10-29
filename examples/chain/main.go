package main

import (
	"context"
	"log"

	. "github.com/twoojoo/ttask/operator"
	. "github.com/twoojoo/ttask/task"
)

func main() {
	chained := T(
		Injectable[string]("t2"),
		Print[string]("> second:"),
	).Catch(func(m *Meta, e error) {
		log.Fatal(e)
	})

	t := T(T(T(
		Injectable[string]("t1"),
		Print[string]("> first:")),
		Chain[string](chained)),
		Print[string]("> third:"),
	).Catch(func(m *Meta, e error) {
		log.Fatal(e)
	})

	err := t.Inject(context.Background(), "message")
	if err != nil {
		log.Fatal(err)
	}
}
