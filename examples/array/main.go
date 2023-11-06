package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	. "github.com/twoojoo/ttask"
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
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Run(ctx)

	fmt.Println("-----------------")

	T(T(
		FromItem("iterate-array-example", []int{0, 1, 2, 3, 4, 5, 5, 6, 7, 8, 9}),
		IterateArray[int]()),
		Print[int]("order >"),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Run(ctx)

	fmt.Println("-----------------")

	T(T(T(T(
		FromItem("parallelize-array-example", []int{0, 1, 2, 3, 4, 5, 5, 6, 7, 8, 9}),
		Distinct[int]()),
		ParallelizeArray[int]()),
		Delay[int](time.Duration(rand.Intn(100))*time.Millisecond)),
		Print[int]("chaos >"),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Run(ctx)
}
