package main

import (
	"context"
	"log"
	"time"

	. "github.com/twoojoo/ttask"
)

func main() {
	ctx := context.Background()

	T(T(T(
		FromStdin("t1"),
		Tap(func(x string) {
			log.Println("> received", x)
		})),
		Delay[string](time.Second)),
		ToStdoutln(func(x string) string {
			return x
		}),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Run(ctx)
}
