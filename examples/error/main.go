package main

import (
	"context"
	"errors"
	"log"

	. "github.com/twoojoo/ttask"
)

func main() {
	count := "1"

	t := T(T(T(T(T(
		Injectable[string]("t1"),
		Print[string]("received >")),
		WithCustomKey[string](func(_ string) string { return "default" })),
		WithContextValue("k1", func(x string) any {
			log.Println("extrancting ctx value... - " + count)
			return x + " (put in ctx) - " + count
		})),
		TapRaw(func(i *Inner, _ *Message[string]) {
			log.Println(i.Context.Value("k1").(string))
		})),
		TapRaw(func(i *Inner, x *Message[string]) {
			err := errors.New("I wanted to throw this error - " + count)
			i.Error(err, "TapRaw:", x.Key, "-")
		}),
	).Catch(func(i *Inner, e error) {
		val := i.Context.Value("k1").(string)
		log.Println("ctx value was:", val)
		log.Printf("[%s] error at %s", i.TaskID(), e.Error())
	}).Lock()

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
