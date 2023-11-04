package main

import (
	"context"
	"log"

	. "github.com/twoojoo/ttask"
)

func main() {
	t := T(T(T(
		Injectable[string]("t1"),
		Print[string]("> 1:")),
		Parallelize[string](5)),
		Print[string]("> 2:"),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Lock()

	for i := 0; i < 6; i++ {
		err := t.Inject(context.Background(), "msg")
		if err != nil {
			log.Fatal(err)
		}
	}
}
