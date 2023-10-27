package task

import (
	"context"
)

// Exec the first step of a Task (and other steps cascading).
// Use this when not manually injecting items in the Task.
func (t *TTask[O, T]) Run(c context.Context) *TTask[O, T] {
	t.run(c)
	return t
}

func (t *TTask[O, T]) run(c context.Context, x ...O) (*T, bool) {
	t.meta.error = nil
	t.meta.Ctx = c

	var msg any
	msg = NewEmptyMessage()
	if len(x) > 0 {
		msg = NewMessage[O](x[0])
	}

	t.meta.ExecNext(msg, t.first)

	if t.meta.error != nil {
		if t.catcher != nil {
			t.catcher(t.meta, t.meta.error)
		}

		return nil, false
	}

	lrMsg, ok := t.meta.lastResult.(*Message[T])
	if !ok {
		panic("not ok msg")
	}

	return &lrMsg.Value, true
}

// Push an item to the Task. Use this when not using a task source.
func (t *TTask[O, T]) Inject(c context.Context, x O) (*T, bool) {
	return t.run(c, x)
}

// Catch any error that was raised in the Task with the m.Error function.
func (t *TTask[O, T]) Catch(catcher func(m *Meta, e error)) *TTask[O, T] {
	t.catcher = catcher
	return t
}