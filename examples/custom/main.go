package main

import (
	"context"
	"log"
	"strconv"

	. "github.com/twoojoo/ttask/operator"
	. "github.com/twoojoo/ttask/task"
)

func main() {

	customOperator := MapRaw[string, int](func(m *Meta, x *Message[string]) int {
		num, err := strconv.Atoi(x.Value)
		if err != nil {
			m.Error(err)
			return 0
		}

		return num
	})


	t := T(T(T(
		Injectable[string]("t1"),
		Print[string]("string >")),
		customOperator),
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
