package operator

import (
	"time"

	"github.com/twoojoo/ttask/task"
)

// Delay the next task step.
func Delay[T any](d time.Duration) task.Operator[T, T] {
	return func(inner *task.Inner, x *task.Message[T], next *task.Step) {
		time.Sleep(d)
		inner.ExecNext(x, next)
	}
}

// Chain another task to the current one syncronously.
// Chaining a locked task will cause the application to panic.
// The act of chaininga locks the chained task as if Lock() method was called.
func Chain[O, T any](t *task.TTask[O, T]) task.Operator[O, T] {
	chainCh := make(chan chainInfo)
	task.T(t, chain[T](chainCh))
	t.Lock()

	return func(inner *task.Inner, x *task.Message[O], next *task.Step) {
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
	innerPtr *task.Inner
	nextPtr  *task.Step
}

func chain[T any](ch chan chainInfo) task.Operator[T, T] {
	return TapRaw(func(inner *task.Inner, x *task.Message[T]) {
		chainInfo := <-ch
		chainInfo.innerPtr.ExecNext(x, chainInfo.nextPtr)
	})
}

// Create an asyncronous branch from the current task using another task.
// When branching, the parent task and the child task will continue their flow concurrently.
// Branching a task will cause the child task to be locked as if Lock() method was called.
// An already locked task can be used as child task when branching.
func Branch[T any](t *task.TTask[T, T]) task.Operator[T, T] {
	t.Lock()

	return func(inner *task.Inner, x *task.Message[T], next *task.Step) {
		msgCopy := *x

		go func() {
			t.InjectRaw(inner.Context, &msgCopy)
		}()

		inner.ExecNext(x, next)
	}
}

// Similar to the Branch operator, but redirects to the new task only messages that pass the provided filter
func BranchWhere[T any](t *task.TTask[T, T], filter func(x T) bool) task.Operator[T, T] {
	t.Lock()

	return func(inner *task.Inner, x *task.Message[T], next *task.Step) {
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
func BranchSwitch[T any](t *task.TTask[T, T], filter func(x T) bool) task.Operator[T, T] {
	t.Lock()

	return func(inner *task.Inner, x *task.Message[T], next *task.Step) {
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
func Parallelize[T any](n int) task.Operator[T, T] {
	cache := []*task.Message[T]{}
	ch := make(chan struct{}, n)

	return func(inner *task.Inner, x *task.Message[T], next *task.Step) {
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
