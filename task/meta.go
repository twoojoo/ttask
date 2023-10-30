package task

import (
	"context"
	"fmt"
	"reflect"
)

// Task metadata and methods to be used inside operators.
type Meta struct {
	Context    context.Context
	taskId     string
	lastResult any
	catcher    func(t *Meta, e error)
	error      error
}

// Calling this function will cause the Task flow to be interrupted before the next operator.
// Use decorator to generate a more detailed error: "[dec1] [dec2] ... err.Error()".
//
// Returining immediatelly after calling this funciton is highly suggested in order to avoid
// unwanted code executions (returned value doesn't matter).
// If the Catch method of the Task isn't used, the error will be lost.
func (m *Meta) Error(e error, decorators ...any) {
	if len(decorators) > 0 {
		dec := ""
		for i, d := range decorators {
			if i != 0 {
				dec += " "
			}

			dec += fmt.Sprint(d)
		}

		m.error = fmt.Errorf("%s %w", dec, e)
	} else {
		m.error = e
	}
}

// Trigger the next Task step.
// Use this in a raw Operator to handle the Task flow in a custom way
// (NOT TYPE SAFE)
func (m *Meta) ExecNext(x any, next *Step) {
	m.lastResult = x

	if m.error != nil {
		if m.catcher != nil {
			m.catcher(m, m.error)
		}

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

func (m Meta) TaskID() string {
	return m.taskId
}
