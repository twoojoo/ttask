package main

import (
	"context"
	"log"
	"strconv"

	. "github.com/twoojoo/ttask/operator"
	. "github.com/twoojoo/ttask/source"
	. "github.com/twoojoo/ttask/task"
)

func main() {
	ctx := context.Background()

	T(T(T(T(T(T(
		FromItem("array-example", [][]int{
			{0, 1, 2, 3, 4},
			{5, 6, 7, 8, 9},
		}),
		FlatArray[int]()),
		MapArray(func(x int) int {
			return x * 3
		})),
		FilterArray(func(x int) bool {
			return x%2 == 0
		})),
		MapArrayRaw(func(i *Inner, x *Message[int]) string {
			return strconv.Itoa(x.Value)
		})),
		ReduceArray("", func(acc *string, x string) string {
			return *acc + x
		})),
		Print[string](">"),
	).Catch(func(m *Inner, e error) {
		log.Fatal(e)
	}).Run(ctx)

}
