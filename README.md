## TTask

A stream processing library for Go, heavily inspired by [alyxstream](https://github.com/smartpricing/alyxstream).

> **⚠** I'm writing this module for learning purposes
>
> **⚠** This module is **not idiomatic Go**. Due to the lack of generics on struct methods (*go 1.21.3* - see [this issue](https://github.com/golang/go/issues/49085)), I was forced to use a weird pattern to replicate a sort of fluent syntax while mantaining full type safety, hence the name of the module. If a future version of go support this feature, a new version of this module may be written.

Working examples in the [examples](./examples) folder.

### Table of contents

- [Importing](#importing)
- [Task](#task)
- [Sources](#sources)
- [Sinks](#sinks)
- [Operators](#operators)
- [Task flow](#task-flow)
- [Windowing](#windowing)
	- [Tumbling window](#tumbling-window)
	- [Hopping window](#hopping-window)
	- [Session window](#session-window)
	- [Counting window](#coutning-window)
- [Raw operators](#raw-operators)
	- [Task meta](#task-meta)
	- [Error handling](#error-handling)
	- [Custom operators](#custom-operators)
- [Next steps](#next-steps)

### Importing

Even if it's not idiomatic, I suggest to import ttask packages in this way, otherwise your code will end up being a bit too verbose (golint will complain).

```go
import (
	. "github.com/twoojoo/ttask"
)
```

### Task

The main abstraction of the library is the Task, which is a set of ordered [operations](#operators) on a stream of data. Tasks can also be [chained and branched](#task-flow) in custom ways. 

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
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Lock()

err := t.Inject(context.Background(), 123)
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
func Map[T, R any](mapper func(x T) R) 

// Filter messages.
func Filter[T, R any](filter func(x T) bool) 

// Perform an action for the message.
func Tap[T any](action func(x T)) 

// Delay the next task step.
func Delay[T any](d time.Duration) 
```

#### Array operators

> for array operators, the generic type refers to the type of the elements of the array

```go
func MapArray[T, R any](mapper func(x T) R) 

// Execute an action for each element of the array
func ForEach[T any](action func(x T)) 

func FilterArray[T any](filter func(x T) bool)

// JS-like array reducer
func ReduceArray[T, R any](base R, reducer func(acc *R, x T) R)

// Flattens an array of type [][]T
func FlatArray[T any]()
```

#### Context operators

```go
// Cache a key/value record in the Task context. Use an extractor function to pull the value from the processed item.
func WithContextValue[T any](k any, ext func(x T) any) 
```

#### Kafka operators

```go
// Perform a commit on the current kafka message.
func KafkaCommit[T any](consumer *kafka.Consumer, logger bool) Operator[types.KafkaMessage[T], types.KafkaMessage[T]]
```

### Windowing

TBD

#### Tumbling window

#### Hopping window

#### Session window

#### Counting window

### Task flow

Chaining/branching tasks is done with some special operators that allow to bisec the task in different ways.

```go
// Chain another task to the current one syncronMapRawously.
// Chaining a locked task will cause the application to panic.
// The act of chaining locks the child task as if Lock() method was called.
// The child task must be an injectable task, otherwise the process will panic.
func Chain[O, T any](t *TTask[O, T]) 

// Create an asyncronous branch from the current task using another task.
// When branching, the parent task and the child task will continue their flow concurrently.
// Branching a task will cause the child task to be locked as if Lock() method was called.
// An already locked task can be used as child task when branching.
// The child task must be an injectable task, otherwise the process will panic.
func Branch[T any](t *TTask[T, T]) 

// Similar to the Branch operator, but redirects to the new task only messages that pass the provided filter
func BranchWhere[T any](t *TTask[T, T], filter func(x T) bool) 

// Similar to the BranchWhere operator, but messages will either pass to the new branch or continue in the current one
func BranchSwitch[T any](t *TTask[T, T], filter func(x T) bool) 

// Process n messages in parallel using an in-memory buffer
func Parallelize[T any](n int)

// Continue the task execution for each element of the array synchronously
func IterateArray[T any]()

// Continue the task exection for each element of the array asynchronously
func ParallelizeArray[T any]()
```

### Raw operators

Most operators have a so called **raw** version, meaning that it give access to lower level task resources, namely *task inner methods and properties* and *message metadata*. For example the Tap operator has its own raw version:

```go
func TapRaw[T any](action func(i *Inner, x *Message[T])) 
```

#### Task meta

TBD

#### Error handling

When the logic of your operator's callback is error prone, it's highly suggested to use the raw version of that operator, since it gives acess to the embedded error handling throgh the Task inner API.

```go
t := T(T(T(
	Injectable[string]("t1"),
	Delay[string](time.Second)),
	MapRaw(func (i *Inner, x Message[string]) int {
		num, err := strconv.Atoi(x)
		if err != nil {
			i.Error(err)
			return 0
		}

		return num
	})),
	Print[int](">"),
).Catch(func(i *Inner, e error) {
	log.Fatal(e)
}).Lock()

t.Inject(context.Background(), "abc")
```

Here's what's happpening:
- since the injected message isn't a numeric string, the Atoi function will return an error
- the error is passed to the **Error** function of the task inner API
- the function returns a 0 (it could be any int value, it will be discarded)
- the error will be catched by the catcher callback passed to the Catch method of the task
- **the task won't continue** to the following operators

> **⚠** if the catcher callback is not set, the error will be raised, but it will be ignored

#### Custom operators

Custom operators (as well as sources and sinks) can be built starting from the **MapRaw** operator.

You can find an example [here](./examples/custom/main.go)

### Next steps

- fill up the readme file with missing stuff (windows, storage, etc..)
- better Memory storage with mutexes
- write a Redis storage for windows operators
- write a Rocksdb storage for windows operators
- write a Cassandra storage for windows operators
- improve persistence/recovery in the task
	- maybe enhancing storage systems
	- using checkpoints (e.g. Checkpoint(Map(func (..) { ...})))
- add WithTimeout and WithDeadline context operators
- add more sources and sinks (e.g. CSV, RabbitMQ, Apache Pulsar, etc..)
