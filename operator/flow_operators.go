package operator

import (
	"time"

	"github.com/twoojoo/ttask/task"
)

// Delay the next task step.
func Delay[T any](d time.Duration) task.Operator[T, T] {
	return func(m *task.Inner, x *task.Message[T], next *task.Step) {
		time.Sleep(d)
		m.ExecNext(x, next)
	}
}

// Chain another task to the current one syncronously.
// Chaining a locked task will cause the application to panic.
// The act of chaininga locks the chained task as if Lock() method was called.
func Chain[O, T any](t *task.TTask[O, T]) task.Operator[O, T] {
	chainCh := make(chan chainInfo)
	task.T(t, chain[T](chainCh))
	t.Lock()

	return func(m *task.Inner, x *task.Message[O], next *task.Step) {
		go func() {
			chainCh <- chainInfo{
				metaPtr: m,
				nextPtr: next,
			}
		}()

		t.InjectRaw(m.Context, x)
	}
}

type chainInfo struct {
	metaPtr *task.Inner
	nextPtr *task.Step
}

func chain[T any](ch chan chainInfo) task.Operator[T, T] {
	return TapRaw(func(m *task.Inner, x *task.Message[T]) {
		chainInfo := <-ch
		chainInfo.metaPtr.ExecNext(x, chainInfo.nextPtr)
	})
}

// Create an asyncronous branch from the current task using another task.
// When branching, the parent task and the child task will continue their flow concurrently.
// Branching a task will cause the child task to be locked as if Lock() method was called.
// An already locked task can be used as child task when branching.
func Branch[T any](t *task.TTask[T, T]) task.Operator[T, T] {
	t.Lock()

	return func(m *task.Inner, x *task.Message[T], next *task.Step) {
		msgCopy := *x

		go func() {
			t.InjectRaw(m.Context, &msgCopy)
		}()

		m.ExecNext(x, next)
	}
}

// Similar to the Branch operator, but redirects to the new task only messages that pass the provided filter
func BranchWhere[T any](t *task.TTask[T, T], filter func(x T) bool) task.Operator[T, T] {
	t.Lock()

	return func(m *task.Inner, x *task.Message[T], next *task.Step) {
		msgCopy := *x

		if filter(x.Value) {
			go func() {
				t.InjectRaw(m.Context, &msgCopy)
			}()
		}

		m.ExecNext(x, next)
	}
}

// Similar to the BranchWhere operator, but messages will either pass to the new branch or continue in the current one
func BranchSwitch[T any](t *task.TTask[T, T], filter func(x T) bool) task.Operator[T, T] {
	t.Lock()

	return func(m *task.Inner, x *task.Message[T], next *task.Step) {
		msgCopy := *x

		if filter(x.Value) {
			go func() {
				t.InjectRaw(m.Context, &msgCopy)
			}()
		} else {

			m.ExecNext(x, next)
		}
	}
}

func Parallelize[T any](n int) task.Operator[T, T] {
	cache := []*task.Message[T]{}
	ch := make(chan int, n)

	return func(m *task.Inner, x *task.Message[T], next *task.Step) {
		cache = append(cache, x)

		if len(cache) == n {
			for _, msg := range cache {
				msgCp := &*msg

				go func() {
					m.ExecNext(msgCp, next)
				}()
			}

			for i := 0; i < n; i++ {
				<-ch
			}

			ch = make(chan int, n)
		}
	}
}
