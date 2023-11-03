package task

import (
	"time"

	"github.com/google/uuid"
)

// Wraps the value that is being processed inside the task. Accessible using the Raw version of operators.
type Message[T any] struct {
	id            string
	EventTime     time.Time
	injestionTime time.Time
	Key           string
	Value         T
}

func NewMessage[T any](value T) *Message[T] {
	now := time.Now()
	return &Message[T]{
		id:            uuid.NewString(),
		injestionTime: now,
		EventTime:     now,
		Value:         value,
	}
}

func NewEmptyMessage() *Message[any] {
	now := time.Now()
	return &Message[any]{
		id:            uuid.NewString(),
		injestionTime: now,
		EventTime:     now,
		Value:         "",
	}
}

func (m Message[T]) GetID() string {
	return m.id
}

func (m Message[T]) GetInjestionTime() time.Time {
	return m.injestionTime
}

func (m *Message[T]) WithKey(k string) *Message[T] {
	m.Key = k
	return m
}

func ReplaceValue[T, R any](m *Message[T], v R) *Message[R] {
	return &Message[R]{
		id:            m.id,
		EventTime:     m.EventTime,
		injestionTime: m.injestionTime,
		Key:           m.Key,
		Value:         v,
	}
}

func ToArray[T any](m *Message[T], elems []Message[T]) *Message[[]T] {
	elemsValues := []T{}
	for _, v := range elems {
		elemsValues = append(elemsValues, v.Value)
	}

	return &Message[[]T]{
		id:            m.id,
		EventTime:     m.EventTime,
		injestionTime: m.injestionTime,
		Key:           m.Key,
		Value:         elemsValues,
	}
}
