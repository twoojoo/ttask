package task

// import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

// Wraps the value that is being processed inside the task. Accessible using the Raw version of operators.
type Message[T any] struct {
	Key string
	Value T
}

func NewMessage[T any](value T) *Message[T] {
	return &Message[T]{
		Value: value,
	}
}

func NewEmptyMessage() *Message[any] {
	return &Message[any]{
		Value: "",
	}
}

func (m *Message[T]) WithKey(k string) *Message[T] {
	m.Key = k
	return m
}

func ReplaceValue[T, R any](m *Message[T], v R) *Message[R] {
	return &Message[R]{
		Key: m.Key,
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
