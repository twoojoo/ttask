package main

import (
	"context"
	"log"

	. "github.com/twoojoo/ttask"
)

func main() {
	ctx := context.Background()

	numbers := [][]int{
		{1, 2, 3, 4},
		{0, 1, 2, 3, 4},
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

	T(T(T(T(
		FromItem("average-example", numbers),
		FlatArray[int]()),
		MapArray(func(x int) float64 {
			return float64(x)
		})),
		Average[float64]()),
		Print[float64]("average >"),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Run(ctx)
}
