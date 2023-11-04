package ttask

import (
	"time"
)

// Delay the next task step.
func Delay[T any](d time.Duration) Operator[T, T] {
	return func(inner *Inner, x *Message[T], next *Step) {
		time.Sleep(d)
		inner.ExecNext(x, next)
	}
}

// Chain another task to the current one syncronously.
// Chaining a locked task will cause the application to panic.
// The act of chaininga locks the chained task as if Lock() method was called.
// The child task must be an injectable task, otherwise the process will panic.
func Chain[O, T any](t *TTask[O, T]) Operator[O, T] {
	if !t.IsInjectable() {
		panic("TTask error: cannot use a non injectable task in  a Branch operator")
	}

	chainCh := make(chan chainInfo)
	Via(t, chain[T](chainCh))
	t.Lock()

	return func(inner *Inner, x *Message[O], next *Step) {
		go func() {
			chainCh <- chainInfo{
				innerPtr: inner,
				nextPtr:  next,
			}
		}()

		t.InjectRaw(inner.Context, x)
	}
}

type chainInfo struct {
	innerPtr *Inner
	nextPtr  *Step
}

func chain[T any](ch chan chainInfo) Operator[T, T] {
	return TapRaw(func(inner *Inner, x *Message[T]) {
		chainInfo := <-ch
		chainInfo.innerPtr.ExecNext(x, chainInfo.nextPtr)
	})
}

// Create an asyncronous branch from the current task using another task.
// When branching, the parent task and the child task will continue their flow concurrently.
// Branching a task will cause the child task to be locked as if Lock() method was called.
// An already locked task can be used as child task when branching.
// The child task must be an injectable task, otherwise the process will panic.
func Branch[T any](t *TTask[T, T]) Operator[T, T] {
	t.Lock()

	if !t.IsInjectable() {
		panic("TTask error: cannot use a non injectable task in  a Branch operator")
	}

	return func(inner *Inner, x *Message[T], next *Step) {
		msgCopy := *x

		go func() {
			t.InjectRaw(inner.Context, &msgCopy)
		}()

		inner.ExecNext(x, next)
	}
}

// Similar to the Branch operator, but redirects to the new task only messages that pass the provided filter
func BranchWhere[T any](t *TTask[T, T], filter func(x T) bool) Operator[T, T] {
	t.Lock()

	if !t.IsInjectable() {
		panic("TTask error: cannot use a non injectable task in  a BranchWhere operator")
	}

	return func(inner *Inner, x *Message[T], next *Step) {
		msgCopy := *x

		if filter(x.Value) {
			go func() {
				t.InjectRaw(inner.Context, &msgCopy)
			}()
		}

		inner.ExecNext(x, next)
	}
}

// Similar to the BranchWhere operator, but messages will either pass to the new branch or continue in the current one
func BranchSwitch[T any](t *TTask[T, T], filter func(x T) bool) Operator[T, T] {
	t.Lock()

	if !t.IsInjectable() {
		panic("TTask error: cannot use a non injectable task in  a BranchSwitch operator")
	}

	return func(inner *Inner, x *Message[T], next *Step) {
		msgCopy := *x

		if filter(x.Value) {
			go func() {
				t.InjectRaw(inner.Context, &msgCopy)
			}()
		} else {
			inner.ExecNext(x, next)
		}
	}
}

// Process n messages in parallel using an in-memory buffer
func Parallelize[T any](n int) Operator[T, T] {
	cache := []*Message[T]{}
	ch := make(chan struct{}, n)

	return func(inner *Inner, x *Message[T], next *Step) {
		cache = append(cache, x)

		if len(cache) == n {
			for _, msg := range cache {
				msgCp := &*msg

				go func() {
					inner.ExecNext(msgCp, next)
					ch <- struct{}{}
				}()
			}

			//wait for all the iterations to complete
			for i := 0; i < n; i++ {
				<-ch
			}

			ch = make(chan struct{}, n)
		}
	}
}
