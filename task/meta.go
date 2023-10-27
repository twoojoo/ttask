package task

import (
	"context"
	"reflect"
)

type Meta struct {
	lastResult any
	Context        context.Context
	error      error
}

// Calling this function will cause the Task flow to be interrupted before the next operator.
// Returining immediatelly after calling this funciton is suggested to avoid unwanted code executions (returned value doesn't matter).
// If the Catch method of the Task isn't used, the error will be lost.
func (m *Meta) Error(e error) {
	m.error = e
}

// Trigger the next Task step. 
// Use this in a raw Operator to handle the Task flow in a custom way 
// (NOT TIPE SAFE)
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
