## TTask

A stream processing library for Go, heavily inspired by [alyxstream](https://github.com/smartpricing/alyxstream).

> **!!!** This module is **not idiomatic Go**. Due to the lack of generics on struct methods (*go 1.21.3*), I was forced to use a weird pattern to replicate a sort of fluent syntax while mantaining full type safety, hence the name of the module. If a future version of go support this feature, a new version of this module may be written.

### Importing

Even if it's not idiomatic, I suggest to import packages in this way, otherwise your code will be a bit verbose.

```go
import (
	. "github.com/twoojoo/ttask/operator"
	. "github.com/twoojoo/ttask/source"
	. "github.com/twoojoo/ttask/task"
	. "github.com/twoojoo/ttask/sink"
	. "github.com/twoojoo/ttask/window"
	. "github.com/twoojoo/ttask/storage"
)
```

### Task

The main abstraction of the library is the Task, wich is a set of ordered operations on a stream of data. Tasks can also be [chained and branched](#chainingandbranching) in custom ways. 

There are two types of task:

- **Injectable** (messages can be pushed programmatically using the **Inject** method)
- based on a autonomous **Source**, e.g. Kafka consumer, Files, etc.. (using the Inject method here will return an error)

```go
t := T(T(T(
		Injectable[int]("t1"),
		Delay[int](time.Second)),
		Map(func (x int) string {
			return strconv.Itoa(x)
		})),
		Print[string](">"),
	).Catch(func(m *Meta, e error) {
		log.Fatal(e)
	})

err := t.Inject(context.Background(), "msg")
```

> NOTE: if the operator can't infer the message type from a given callback, the type must be provided as generic to the operator itself (e.g. Print and Delay operators) 

### Source

### Sink

### Operators

### Windows

### Chaining and branching