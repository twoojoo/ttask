package main

import (
	"context"
	"log"

	. "github.com/twoojoo/ttask/operator"
	. "github.com/twoojoo/ttask/storage"
	. "github.com/twoojoo/ttask/source"
	. "github.com/twoojoo/ttask/task"
	. "github.com/twoojoo/ttask/window"
)

func main() {
	mem := &MemoryStorage[Message[string]]{}

	T(T(FromStringSplit("ciao mi chiamo Giovanni e sono della provincia di Treviso", " "),
		TumblingWindowCount[string](mem, 2)),
		Tap((func(x []string) { log.Println(x) })),
	).Catch(func(m *Meta, e error) {
		log.Fatal(e)
	}).Run(context.Background())
}
