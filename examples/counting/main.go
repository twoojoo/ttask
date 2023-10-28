package main

import (
	"context"
	"log"
	"time"

	. "github.com/twoojoo/ttask/operator"
	. "github.com/twoojoo/ttask/source"
	. "github.com/twoojoo/ttask/storage"
	. "github.com/twoojoo/ttask/task"
	. "github.com/twoojoo/ttask/window"
)

func main() {

	T(T(
		FromStringSplit("one two three four five six seven eight nine ten", " "),
		CountingWindow(CWOptions[string]{
			Storage: Memory[string](),
			Id:      "win1",
			Size:    2,
		})),
		Print[[]string](">"),
	).Catch(func(m *Meta, e error) {
		log.Fatal(e)
	}).Run(context.Background())


	T(T(
		FromInterval(time.Second, 10, func(count int) int { return count }),
		CountingWindow(CWOptions[int]{
			Storage:       Memory[int](),
			Id:            "win2",
			Size:          2,
			MaxInactivity: 1500 * time.Millisecond,
		})),
		Print[[]int](">"),
	).Catch(func(m *Meta, e error) {
		log.Fatal(e)
	}).Run(context.Background())

}

// Memory[int](), "win2", 2
