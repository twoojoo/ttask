package main

import (
	"context"
	"log"
	"strconv"

	. "github.com/twoojoo/ttask/operator"
	. "github.com/twoojoo/ttask/task"
)

// generate a custom operator that use custom params
func customOperator(toSum int) Operator[string, int] {
	return MapRaw[string, int](func(m *Meta, x *Message[string]) int {
		// opertator logic start:

		num, err := strconv.Atoi(x.Value)
		
		// graceful error handling
		if err != nil {
			m.Error(err)
			return 0
		}

		// next step message value
		return num + toSum
	
		// operator logic end.
	})
}

func main() {

	t := T(T(T(
		Injectable[string]("t1"),
		Print[string]("string >")),
		customOperator(2)),
		Print[int]("integer >"),
	).Catch(func(m *Meta, e error) {
		val := m.Context.Value("k1").(string)
		log.Println("ctx value was:", val)
		log.Println("ERROR:", e)
	})

	err := t.Inject(context.Background(), "1")
	if err != nil {
		log.Fatal(err)
	}
}
