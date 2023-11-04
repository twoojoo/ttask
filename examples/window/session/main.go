package main

import (
	"context"
	"log"
	"time"

	. "github.com/twoojoo/ttask"
)

func main() {
	T(T(
		FromInterval("t1", time.Second, 10, func(count int) int { return count }),
		SessionWindow("win-1", SWOptions[int]{MaxInactivity: 700 * time.Millisecond})),
		Print[[]int]("1 at time >"),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Run(context.Background())

	time.Sleep(2 * time.Second)

	T(T(
		FromInterval("t1", time.Second, 10, func(count int) int { return count }),
		SessionWindow("win-1", SWOptions[int]{
			MaxSize:       1500 * time.Millisecond,
			MaxInactivity: 1010 * time.Millisecond,
		})),
		Print[[]int]("2 at time >"),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Run(context.Background())

	time.Sleep(2 * time.Second)
}
