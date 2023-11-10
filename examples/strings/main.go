package main

import (
	"context"
	"log"
	"time"

	. "github.com/twoojoo/ttask"
)

func main() {
	ctx := context.Background()

	str := "i topi non avevano nipoti"

	T(T(T(
		FromStringSplit("strings-example", str, " "),
		CountingWindow("win1", CWOptions[string]{
			MaxInactivity: time.Second,
		})),
		Concat("|")),
		Print[string]("sum >"),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Run(ctx)
}
