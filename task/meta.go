package task

import (
	"context"
	"reflect"
)

type Meta struct {
	lastResult any
	Ctx        context.Context
	error      error
}

func (t *Meta) Error(e error) {
	t.error = e
}

func (t *Meta) ContextValue(k any) any {
	return t.Ctx.Value(k)
}

func (m *Meta) ExecNext(x any, next *Step) {
	m.lastResult = x

	if m.error != nil {
		return
	}

	if next == nil {
		return
	}

	nextActionValue := reflect.ValueOf(next.action)

	if nextActionValue.Kind() != reflect.Func {
		panic("operator is not a func!")
	}

	argsValue := []reflect.Value{
		reflect.ValueOf(m),
		reflect.ValueOf(x),
		reflect.ValueOf(next.next),
	}

	nextActionValue.Call(argsValue)
}
