package main

import (
	"context"
	"log"
	"time"

	. "github.com/twoojoo/ttask/operator"
	. "github.com/twoojoo/ttask/source"
	. "github.com/twoojoo/ttask/task"
	. "github.com/twoojoo/ttask/window"
)

func main() {
	T(T(
		FromInterval("thw", time.Second, 10, func(count int) int { return count }),
		HoppingWindow(HWOptions[int]{
			Size: 1900 * time.Millisecond,
			Hop: 2 * time.Second,
		})),
		Print[[]int](">"),
	).Catch(func(m *Meta, e error) {
		log.Fatal(e)
	}).Run(context.Background())
}
