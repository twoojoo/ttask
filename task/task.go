package task

type TTask[O, T any] struct {
	id         string
	injectable bool
	first      *Step
	last       int
	path       map[int]any
	meta       *Meta
}

// Use this to build custom sources only. Not an injectable task.
func Task[T any](id string) *TTask[T, T] {
	t := TTask[T, T]{
		id:   id,
		last: 0,
		path: map[int]any{},
		first: &Step{
			action: nil,
			next:   nil,
		},
		meta: &Meta{
			Context: nil,
			error:   nil,
		},
	}

	return &t
}

// Initialize an injectable Task with the first step message type as generic.
// To push messages to this Task use the Inject method.
func Injectable[T any](id string) *TTask[T, T] {
	t := TTask[T, T]{
		id:         id,
		injectable: true,
		last:       0,
		path:       map[int]any{},
		first: &Step{
			action: nil,
			next:   nil,
		},
		meta: &Meta{
			Context: nil,
			error:   nil,
		},
	}

	return &t
}

// Add an operator to the Task. Returns the updated Task.
func T[O, T, R any](t *TTask[O, T], operator Operator[T, R]) *TTask[O, R] {
	if t.last == 0 {
		t.first = &Step{
			action: operator,
			next:   nil,
		}
	} else {
		curr := t.first
		for i := 0; i <= t.last; i++ {
			if curr.next == nil {
				curr.next = &Step{
					action: operator,
					next:   nil,
				}

				break
			}

			curr = curr.next
		}
	}

	return &TTask[O, R]{
		injectable: t.injectable,
		first:      t.first,
		path:       t.path,
		last:       t.last + 1,
		meta:       t.meta,
	}
}

type Operator[T, R any] func(t *Meta, x *Message[T], next *Step)

type Step struct {
	next   *Step
	action any
}
