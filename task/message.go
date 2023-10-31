package task

import "time"

// Wraps the value that is being processed inside the task. Accessible using the Raw version of operators.
type Message[T any] struct {
	EventTime      time.Time
	processingTime time.Time
	Key            string
	Value          T
}

func NewMessage[T any](value T) *Message[T] {
	now := time.Now()
	return &Message[T]{
		processingTime: now,
		EventTime:      now,
		Value:          value,
	}
}

func NewEmptyMessage() *Message[any] {
	now := time.Now()
	return &Message[any]{
		processingTime: now,
		EventTime:      now,
		Value:          "",
	}
}

func (m *Message[T]) WithKey(k string) *Message[T] {
	m.Key = k
	return m
}

func ReplaceValue[T, R any](m *Message[T], v R) *Message[R] {
	return &Message[R]{
		Key:   m.Key,
		Value: v,
	}
}

func ToArray[T any](m *Message[T], elems []Message[T]) *Message[[]T] {
	elemsValues := []T{}
	for _, v := range elems {
		elemsValues = append(elemsValues, v.Value)
	}

	return &Message[[]T]{
		Key:   m.Key,
		Value: elemsValues,
	}
}
