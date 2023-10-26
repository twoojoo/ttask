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

	t := T(T(T(T(T(Task[string](),
		Print[string]("received >")),
		WithContextValue("k1", func(x string) any {
			log.Println("extrancting ctx value...")
			return x + " (put in ctx) " + count
		})),
		Tap(func(x string) {
			log.Println("sjdfnksdjfnk")
		})),
		TapRaw(func(m *Meta, _ *Message[string]) {
			log.Println(m.ContextValue("k1").(string))
		})),
		TapRaw(func(m *Meta, _ *Message[string]) {
			m.Error(errors.New("I wanted to throw this error"))
		})).
		Catch(func(m *Meta, e error) {
			v := m.ContextValue("k1").(string)
			log.Println("ctx value was:", v)
			log.Println(e)
		})

	_, ok := t.Inject(context.Background(), "message 1")
	if !ok {
		log.Println("an error occurred or value task stopped")
	}

	count = "2"
	_, ok = t.Inject(context.Background(), "message 2")
	if !ok {
		log.Println("an error occurred or value task stopped")
	}
}
