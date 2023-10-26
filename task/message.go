package task

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

//Wraps the value that is being processed inside the task. Accessible using the Raw version of operators.
type Message[T any] struct {
	Key            string
	TopicPartition []kafka.TopicPartition
	Value          T
}

func NewMessage[T any](value T) *Message[T] {
	return &Message[T]{
		Value:          value,
		TopicPartition: []kafka.TopicPartition{},
	}
}

func NewEmptyMessage() *Message[any] {
	return &Message[any]{
		Value:          "",
		TopicPartition: []kafka.TopicPartition{},
	}
}

func (m *Message[T]) WithTopicPartition(tp kafka.TopicPartition) *Message[T] {
	m.TopicPartition = append(m.TopicPartition, tp)
	return m
}

func (m *Message[T]) WithKey(k string) *Message[T] {
	m.Key = k
	return m
}

func ReplaceValue[T, R any](m *Message[T], v R) *Message[R] {
	return &Message[R]{
		Key: m.Key,
		TopicPartition: m.TopicPartition,
		Value: v,
	}
}
