package ttask

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Wraps the value that is being processed inside the task. Accessible using the Raw version of operators.
type Message[T any] struct {
	Id            string
	EventTime     time.Time
	InjestionTime time.Time
	Key           string
	Value         T
}

func NewMessage[T any](value T) *Message[T] {
	now := time.Now()
	return &Message[T]{
		Id:            uuid.NewString(),
		InjestionTime: now,
		EventTime:     now,
		Value:         value,
	}
}

func newEmptyMessage() *Message[any] {
	now := time.Now()
	return &Message[any]{
		Id:            uuid.NewString(),
		InjestionTime: now,
		EventTime:     now,
		Value:         "",
	}
}

func (m Message[T]) GetID() string {
	return m.Id
}

func (m Message[T]) GetInjestionTime() time.Time {
	return m.InjestionTime
}

func (m *Message[T]) withKey(k string) *Message[T] {
	m.Key = k
	return m
}

func replaceValue[T, R any](m *Message[T], v R) *Message[R] {
	return &Message[R]{
		Id:            m.Id,
		EventTime:     m.EventTime,
		InjestionTime: m.InjestionTime,
		Key:           m.Key,
		Value:         v,
	}
}

func toArray[T any](m *Message[T], elems []Message[T]) *Message[[]T] {
	elemsValues := []T{}
	for _, v := range elems {
		elemsValues = append(elemsValues, v.Value)
	}

	return &Message[[]T]{
		Id:            m.Id,
		EventTime:     m.EventTime,
		InjestionTime: m.InjestionTime,
		Key:           m.Key,
		Value:         elemsValues,
	}
}

func (m *Message[T]) messageToBytes() ([]byte, error) {
	b, err := json.Marshal(*m)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func bytesToMessage[T any](b *[]byte) (*Message[T], error) {
	msg := &Message[T]{}

	err := json.Unmarshal(*b, msg)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return msg, nil
}
