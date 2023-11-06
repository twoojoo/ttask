package main

import (
	"context"
	"log"
	"time"

	. "github.com/twoojoo/ttask"
)

func main() {
	ctx := context.Background()

	T(T(T(T(
		FromStdin("t1"),
		Tap(func(x string) {
			log.Println("> received", x)
		})),
		CountingWindow("w1", CWOptions[string]{
			MaxInactivity: 2 * time.Second,
			Size:          3,
		})),
		ReduceArray("", func(acc *string, x string) string {
			log.Println("> reducing", x)
			return *acc + "|" + x
		})),
		ToStdoutln(func(x string) string {
			return x
		}),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Run(ctx)
}
