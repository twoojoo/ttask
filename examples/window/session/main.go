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
		FromInterval("t1", time.Second, 10, func(count int) int { return count }),
		SessionWindow(SWOptions[int]{MaxInactivity: 1100 * time.Millisecond})),
		Print[[]int]("> window closed"),
	).Catch(func(m *Meta, e error) {
		log.Fatal(e)
	}).Run(context.Background())

	time.Sleep(5 * time.Second)
}