package main

import (
	"context"
	"log"
	"time"

	. "github.com/twoojoo/ttask/operator"
	. "github.com/twoojoo/ttask/task"
)

func main() {
	branch := T(T(
		Injectable[string]("t2"),
		Delay[string](time.Second)),
		Print[string]("> third:"),
	).Catch(func(m *Meta, e error) {
		log.Fatal(e)
	}).Lock()

	t := T(T(T(
		Injectable[string]("t1"),
		Print[string]("> first:")),
		Branch[string](branch)),
		Print[string]("> second:"),
	).Catch(func(m *Meta, e error) {
		log.Fatal(e)
	})

	err := t.Inject(context.Background(), "message")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(2 * time.Second)
}
