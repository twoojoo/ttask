## TTask

A stream processing library for Go, heavily inspired by [alyxstream](https://github.com/smartpricing/alyxstream).

> I'm writing this module for learning purposes

> **!!!** This module is **not idiomatic Go**. Due to the lack of generics on struct methods (*go 1.21.3* - see [this issue](https://github.com/golang/go/issues/49085)), I was forced to use a weird pattern to replicate a sort of fluent syntax while mantaining full type safety, hence the name of the module. If a future version of go support this feature, a new version of this module may be written.

### Table of contents

- [Importing](#importing)
- [Task](#task)
- [Sources](#sources)
- [Sinks](#sinks)
- [Operators](#operators)
- [Windowing](#windowing)
- [Raw operators](#raw-operators)
	- [Task meta](#task-meta)
	- [Error handling](#error-handling)

### Importing

Even if it's not idiomatic, I suggest to import ttask packages in this way, otherwise your code will end up being a bit too verbose.

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

The main abstraction of the library is the Task, wich is a set of ordered operations on a stream of data. Tasks can also be [chained and branched](#chaining-and-branching) in custom ways. 

There are two types of task:

- **Injectable** (messages can be pushed programmatically using the **Inject** method)
- based on an autonomous **Source**, e.g. Kafka consumer, Files, etc.. (using the Inject method here will return an error)

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
	}).Lock()

err := t.Inject(context.Background(), "msg")
```

> NOTE: if the operator can't infer the message type from a given callback (e.g. Map operator), the type must be provided as generic to the operator itself (e.g. Print and Delay operators) 

> **Lock** method prevents the task from being further extended with more operators. Trying to extend a locked task will cause the application to panic.

> **Catch** method allows graceful [error handling](#error-handling).

### Sources

### Sinks

### Operators

### Windowing

### Chaining and branching

### Raw Operators

#### Task meta

#### Error handling

