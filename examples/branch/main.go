package main

import (
	"context"
	"log"
	"time"

	. "github.com/twoojoo/ttask"
)

func main() {
	branch := T(T(
		Injectable[string]("t2"),
		Delay[string](time.Second)),
		Print[string]("> third:"),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Lock()

	// ALTERNATIVE SINTAX

	/**
	 * pay attention:
	 * all variables are pointers to the
	 * same task instance!!
	 * */

	source := Injectable[string]("t1")

	step1 := T(source, Print[string]("> first:"))

	step2 := T(step1, Branch[string](branch))

	task := T(step2, Print[string]("> second:"))

	task.Lock().Catch(func(i *Inner, e error) {
		log.Fatal(e)
	})

	err := task.Inject(context.Background(), "message")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(2 * time.Second)
}
