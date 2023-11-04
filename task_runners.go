package ttask

import (
	"context"
	"errors"
)

// Exec the first step of a Task (and other steps cascading).
// Use this when not manually injecting items in the Task.
// This method also lock the task
func (t *TTask[O, T]) Run(c context.Context) error {
	if t.injectable {
		return errors.New("TTask Error: can't run an injectable task")
	}

	t.Lock()
	t.run(c)

	return nil
}

func (t *TTask[O, T]) run(c context.Context, x ...O) {
	t.inner.error = nil
	t.inner.Context = c

	var msg any
	msg = newEmptyMessage()
	if len(x) > 0 {
		msg = NewMessage[O](x[0])
	}

	t.inner.ExecNext(msg, t.first)
}

func (t *TTask[O, T]) runRaw(c context.Context, x *Message[O]) {
	t.inner.error = nil
	t.inner.Context = c
	t.inner.ExecNext(x, t.first)
}

// Push an item to the Task. Use this when not using a task source.
func (t *TTask[O, T]) Inject(c context.Context, x O) error {
	if !t.injectable {
		return errors.New("TTask Error: can't inject a message in a non-injectable task")
	}

	t.run(c, x)
	return nil
}

func (t *TTask[O, T]) InjectRaw(c context.Context, m *Message[O]) error {
	if !t.injectable {
		return errors.New("TTask Error: can't inject a message in a non-injectable task")
	}

	t.runRaw(c, m)
	return nil
}

// Catch any error that was raised in the Task with the m.Error function.
func (t *TTask[O, T]) Catch(catcher func(i *Inner, e error)) *TTask[O, T] {
	t.inner.catcher = catcher
	return t
}

// Lock the task to prevent it from being further extended with more operators.
func (t *TTask[O, T]) Lock() *TTask[O, T] {
	t.locked = true
	return t
}
