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
	- [Tumbling window](#tumbling-window)
	- [Hopping window](#hopping-window)
	- [Session window](#session-window)
	- [Counting window](#coutning-window)
- [Raw operators](#raw-operators)
	- [Task meta](#task-meta)
	- [Error handling](#error-handling)
	- [Custom operators](#custom-operators)

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

The main abstraction of the library is the Task, wich is a set of ordered [operations](#operators) on a stream of data. Tasks can also be [chained and branched](#chaining-and-branching) in custom ways. 

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

TBD

### Sinks

TBD

### Operators

Sequencial operations on messages that pass through a task are defined by the operators of that task. 

#### Base operators

```go
// Set a custom message key, possibily from the message itself.
func WithCustomKey[T any](extractor func(x T) string) 

// Set a custom message event time, possibly from the message itself.
func WithEventTime[T any](extractor func(x T) time.Time) 

// Print message value with a given prefix.
func Print[T any](prefix ...string) 

// Map the message value.
func Map[T, R any](cb func(x T) R) 

// Filter messages.
func Filter[T, R any](cb func(x T) bool) 

// Perform an action for the message.
func Tap[T any](cb func(x T)) 

// Delay the next task step.
func Delay[T any](d time.Duration) 
```

#### Array operators

```go
func MapArray[T, R any](cb func(x T) R) 

func Each[T any](cb func(x T)) 

func FilterArray[T any](cb func(x T) bool)
```

#### Context operators

```go
//Cache a key/value record in the Task context. Use an extractor function to pull the value from the processed item.
func WithContextValue[T any](k any, ext func(x T) any) 
```

#### Kafka operators

```go
//Perform a commit on the current kafka message.
func KafkaCommit[T any](consumer *kafka.Consumer, logger bool) Operator[types.KafkaMessage[T], types.KafkaMessage[T]]
```

### Windowing

TBD

#### Tumbling window

#### Hopping window

#### Session window

#### Counting window

### Chaining and branching

Chaining/branching tasks is done with some special operators that allow to bisec the task in different ways.

```go
// Chain another task to the current one syncronously.
// Chaining a locked task will cause the application to panic.
// The act of chaininga locks the chained task as if Lock() method was called.
func Chain[O, T any](t *TTask[O, T]) 

// Create an asyncronous branch from the current task using another task.
// When branching, the parent task and the child task will continue their flow concurrently.
// Branching a task will cause the child task to be locked as if Lock() method was called.
// An already locked task can be used as child task when branching.
func Branch[T any](t *TTask[T, T]) 

//Similar to the Branch operator, but redirects to the new task only messages that pass the provided filter
func BranchWhere[T any](t *TTask[T, T], filter func(x T) bool) 

//Similar to the BranchWhere operator, but messages will either pass to the new branch or continue in the current one
func BranchSwitch[T any](t *TTask[T, T], filter func(x T) bool) 
```

### Raw Operators

Most operators have a so called **raw** version, meaning that it give access to lower level task resources, namely *task inner methods and properties* and *message metadata*. For example the Tap operator has its own raw version:

```go
func TapRaw[T any](cb func(m *Meta, x *Message[T])) 
```

#### Task meta

TBD

#### Error handling

When the logic of your operator's callback is error prone, it's highly suggested to use the raw version of that operator, since it gives acess to the enbedded error handling throgh the Task inner API.

```go
t := T(T(T(
	Injectable[string]("t1"),
	Delay[string](time.Second)),
	TapRaw(func (m *Meta, x Message[string]) int {
		num, err := strconv.Atoi(x)
		if err != nil {
			m.Error(err)
			return 0
		}

		return num
	})),
	Print[int](">"),
).Catch(func(m *Meta, e error) {
	log.Fatal(e)
}).Lock()

t.Inject(context.Background(), "abc")
```

Here's what's happpening:
- since the injected message isn't a numeric string, the Atoi function will return an error
- the error is passet to the Error function of the task APi
- the function returns a 0 (it could be any int value, it will be discarded)
- the error will be catched by the catcher callback passed to the Catch method of the task
- in this case the process will stop

> if the catcher callback is not set, the error will be raised, but it will be ignored

#### Custom operators

TBD
