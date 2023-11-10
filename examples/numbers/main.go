package main

import (
	"context"
	"log"

	. "github.com/twoojoo/ttask"
)

func main() {
	ctx := context.Background()

	T(T(T(
		FromItem("numbers-example", [][]int{
			{0, 1, 2, 3, 4},
			{0, 1, 2, 3, 4},
			{15, 16, 17, 18, 19},
			{5, 6, 7, 8, 9},
			{10, 11, 12, 13, 14},
		}),
		FlatArray[int]()),
		Sum[int]()),
		Print[int]("sum >"),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Run(ctx)

	T(T(T(
		FromItem("numbers-example", [][]int{
			{1, 2, 7, 23, 11, 3, 15, 4, 12, 5},
			{15, 16, 17, 18, 19},
			{5, 6, 7, 8, 9},
			{10, 11, 12, 13, 14},
		}),
		FlatArray[int]()),
		Multiply[int]()),
		Print[int]("product >"),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Run(ctx)
}
