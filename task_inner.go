package ttask

import (
	"context"
	"fmt"
	"reflect"
)

// Task metadata and methods to be used inside operators.
type Inner struct {
	Context    context.Context
	taskId     string
	lastResult any
	catcher    func(t *Inner, e error)
	error      error
	storage    Storage
}

func (i *Inner) getStorage() Storage {
	return i.storage
}

func (i *Inner) setStorage(s Storage) {
	i.storage = s
}

// Calling this function will cause the Task flow to be interrupted before the next operator.
// Use decorator to generate a more detailed error: "[dec1] [dec2] ... err.Error()".
//
// Returining immediatelly after calling this funciton is highly suggested in order to avoid
// unwanted code executions (returned value doesn't matter).
// If the Catch method of the Task isn't used, the error will be lost.
func (inner *Inner) Error(e error, decorators ...any) {
	if len(decorators) > 0 {
		dec := ""
		for i, d := range decorators {
			if i != 0 {
				dec += " "
			}

			dec += fmt.Sprint(d)
		}

		inner.error = fmt.Errorf("%s %w", dec, e)
	} else {
		inner.error = e
	}
}

// Trigger the next Task step.
// Use this in a raw Operator to handle the Task flow in a custom way
// (NOT TYPE SAFE)
func (inner *Inner) ExecNext(x any, next *Step) {
	inner.lastResult = x

	if inner.error != nil {
		if inner.catcher != nil {
			inner.catcher(inner, inner.error)
		}

		// msgId := getMessageId(x)
		// inner.storage.clearCheckpoint(inner.TaskID(), msgId)

		return
	}

	if next == nil {
		// msgId := getMessageId(x)
		// inner.storage.clearCheckpoint(inner.TaskID(), msgId)
		return
	}

	nextActionValue := reflect.ValueOf(next.action)

	if nextActionValue.Kind() != reflect.Func {
		panic("operator is not a func!")
	}

	argsValue := []reflect.Value{
		reflect.ValueOf(inner),
		reflect.ValueOf(x),
		reflect.ValueOf(next.next),
	}

	nextActionValue.Call(argsValue)
}

// Return the id of the task
func (m Inner) TaskID() string {
	return m.taskId
}

func getMessageId(msg any) string {
	msgPtr := reflect.ValueOf(msg)

	for i := 0; i < msgPtr.Elem().NumField(); i++ {
		str, ok := msgPtr.Elem().
			MethodByName("GetID").
			Call([]reflect.Value{})[0].
			Interface().(string)

		if !ok {
			panic("failed to extract message id")
		}

		return str
	}

	panic("failed to extract message id")
}
