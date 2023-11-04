package main

import (
	"context"
	"log"
	"time"

	. "github.com/twoojoo/ttask"
)

func main() {
	T(T(T(
		FromInterval("thw", 1000*time.Millisecond, 100, func(count int) int { return count }),
		WithCustomKey(func(x int) string { return "default" })),
		HoppingWindow("win-1", HWOptions[int]{
			Size: 2100 * time.Millisecond,
			Hop:  600 * time.Millisecond,
		})),
		Print[[]int](">"),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Run(context.Background())

	time.Sleep(2 * time.Second)
}

/**
 *
 * 0.........1.........2.........3.........4.........5.........6.........7.........8.........9
 * ----------------------
 *       ----------------------
 *             ----------------------
 *                   ----------------------
 *                         ----------------------
 *                               ----------------------
 *                                     ----------------------
 *                                           ----------------------
 *                                                 ----------------------
 *                                                       ----------------------
 *                                                             ----------------------
 *                                                                   ----------------------
 *                                                                         ----------------------
 *                                                                               ----------------------
 *                                                                                     ----------------------
 *
 *  */
