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
		FromInterval(time.Second, 10, func(count int) int { return count }),
		TumblingWindow(TWOptions[int]{Size: 1500 * time.Millisecond})),
		Print[[]int](">"),
	).Catch(func(m *Meta, e error) {
		log.Fatal(e)
	}).Run(context.Background())

	time.Sleep(10 * time.Second)
}
