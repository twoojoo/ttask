package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	. "github.com/twoojoo/ttask"
)

func main() {
	ctx := context.Background()

	numbers := [][]int{
		{1, 2, 3, 4},
		{15, 16, 17, 18, 19},
		{5, 6, 7, 8, 9},
		{10, 11, 12, 13, 14},
	}

	T(T(T(
		FromItem("sum-example", numbers),
		FlatArray[int]()),
		Sum[int]()),
		Print[int]("sum >"),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Run(ctx)

	T(T(T(
		FromItem("multiply-example", numbers),
		FlatArray[int]()),
		Multiply[int]()),
		Print[int]("product >"),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Run(ctx)

	T(T(T(T(T(
		FromItem("average-example", numbers),
		FlatArray[int]()),
		Array(Filter(IsEven[int]()))),
		Array(ToFloat32[int]())),
		Average[float32]()),
		Print[float32]("average >"),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Run(ctx)

	T(T(
		FromItem("average-array-example", numbers),
		Array(Average[int]())),
		Print[[]int]("array average >"),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Run(ctx)

	/**
	 * - generates a random (0-300) number every 200ms
	 * - collects them in a slice every 3 numbers
	 * - calculates the average of every array
	 * */
	T(T(T(
		FromInterval(
			"window-array-example",
			200*time.Millisecond, 10,
			func(x int) int { return rand.Intn(300) },
		),
		CountingWindow("win1", CWOptions[int]{
			MaxInactivity: time.Second,
			Size:          3,
		})),
		Average[int]()),
		Print[int]("window average >"),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Run(ctx)
}
