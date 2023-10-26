package task

import (
	"reflect"
)

func ExecNext(m *Meta, x any, next *Step) {
	// m.lastResult =  reflect.ValueOf(x).Interface()

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
