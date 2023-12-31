package main

import (
	"context"
	"log"
	"strconv"

	. "github.com/twoojoo/ttask"
)

// generate a custom operator that use custom params
func customOperator(toSum int) Operator[string, int] {
	return MapRaw[string, int](func(inner *Inner, x *Message[string]) int {
		// opertator logic start:

		num, err := strconv.Atoi(x.Value)

		// graceful error handling
		if err != nil {
			inner.Error(err)
			return 0
		}

		// next step message value
		return num + toSum

		// operator logic end.
	})
}

// generate a custom source
func customSource(taskId string, end int) *TTask[any, string] {
	return T(Task[any](taskId), customSourceLogic(end))
}

// define custom source logic
func customSourceLogic(end int) Operator[any, string] {
	return func(inner *Inner, _ *Message[any], next *Step) {
		for i := 0; i < end; i++ {
			num := i * i
			val := strconv.Itoa(num)
			//trigger a task execution (not type safe)
			inner.ExecNext(NewMessage(val), next)
		}
	}
}

func main() {
	T(T(T(
		customSource("t1", 3),
		Print[string]("string >")),
		customOperator(2)),
		Print[int]("integer >"),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Run(context.Background())
}
