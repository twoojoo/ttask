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

	// should split 
	T(T(
		FromStringSplit("t1", "one two three four five six seven eight nine ten", " "),
		CountingWindow(CWOptions[string]{
			Size: 3,
			MaxInactivity: 1000 * time.Millisecond,
		})),
		Print[[]string](">"),
	).Catch(func(m *Meta, e error) {
		log.Fatal(e)
	}).Run(context.Background())

	// time.Sleep(2 * time.Second)


	T(T(
		FromInterval("t2", time.Second, 10, func(count int) int { return count }),
		CountingWindow(CWOptions[int]{
			Size:          2,
			MaxInactivity: 700 * time.Millisecond,
		})),
		Print[[]int](">"),
	).Catch(func(m *Meta, e error) {
		log.Fatal(e)
	}).Run(context.Background())

	T(T(
		FromInterval("t3", time.Second, 10, func(count int) int { return count }),
		CountingWindow(CWOptions[int]{
			Size:          2,
			MaxInactivity: 1010 * time.Millisecond,
		})),
		Print[[]int](">"),
	).Catch(func(m *Meta, e error) {
		log.Fatal(e)
	}).Run(context.Background())

}
