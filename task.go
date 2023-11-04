package ttask

type TTask[O, T any] struct {
	id         string
	injectable bool
	locked     bool
	first      *Step
	last       int
	path       map[int]any
	inner      *Inner
}

// Use this to build custom sources only. Not an injectable task.
func Task[T any](id string) *TTask[T, T] {
	t := TTask[T, T]{
		id:         id,
		last:       0,
		injectable: false,
		locked:     false,
		path:       map[int]any{},
		first: &Step{
			action: nil,
			next:   nil,
		},
		inner: &Inner{
			taskId:  id,
			Context: nil,
			error:   nil,
			storage: NewMemoryStorage(),
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
		locked:     false,
		last:       0,
		path:       map[int]any{},
		first: &Step{
			action: nil,
			next:   nil,
		},
		inner: &Inner{
			taskId:  id,
			Context: nil,
			error:   nil,
		},
	}

	return &t
}

// Add an operator to the Task. Returns the updated Task.
func T[O, T, R any](t *TTask[O, T], operator Operator[T, R]) *TTask[O, R] {
	return Via(t, operator)
}

// Add an operator to the Task. Returns the updated Task.
func Via[O, T, R any](t *TTask[O, T], operator Operator[T, R]) *TTask[O, R] {
	if t.locked {
		panic("adding operator to locked task: " + t.id)
	}

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
		id:         t.id,
		injectable: t.injectable,
		locked:     t.locked,
		first:      t.first,
		path:       t.path,
		last:       t.last + 1,
		inner:      t.inner,
	}
}

func (t *TTask[O, T]) WithStorage(s Storage) *TTask[O, T] {
	t.inner.setStorage(s)
	return t
}

func (t TTask[O, T]) IsInjectable() bool {
	return t.injectable
}

type Operator[T, R any] func(t *Inner, x *Message[T], next *Step)

type Step struct {
	next   *Step
	action any
}
