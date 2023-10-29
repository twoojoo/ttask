package main

import (
	"context"
	"errors"
	"log"

	. "github.com/twoojoo/ttask/operator"
	. "github.com/twoojoo/ttask/task"
)

func main() {
	count := "1"

	t := T(T(T(T(
		Injectable[string]("t1"),
		Print[string]("received >")),
		WithContextValue("k1", func(x string) any {
			log.Println("extrancting ctx value... - " + count)
			return x + " (put in ctx) - " + count
		})),
		TapRaw(func(m *Meta, _ *Message[string]) {
			log.Println(m.Context.Value("k1").(string))
		})),
		TapRaw(func(m *Meta, _ *Message[string]) {
			m.Error(errors.New("I wanted to throw this error - " + count))
		}),
	).Catch(func(m *Meta, e error) {
		val := m.Context.Value("k1").(string)
		log.Println("ctx value was:", val)
		log.Println("ERROR:", e)
	})

	err := t.Inject(context.Background(), "message 1")
	if err != nil {
		log.Fatal(err)
	}

	count = "2"
	err = t.Inject(context.Background(), "message 2")
	if err != nil {
		log.Fatal(err)
	}
}
