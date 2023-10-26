package task

import "context"

type Meta struct {
	lastResult Message[any]
	Ctx        context.Context
	error      error
}

func (t *Meta) Error(e error) {
	t.error = e
}

func (t *Meta) ContextValue(k any) any {
	return t.Ctx.Value(k)
}
