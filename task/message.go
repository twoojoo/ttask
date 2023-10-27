package task

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

// Wraps the value that is being processed inside the task. Accessible using the Raw version of operators.
type Message[T any] struct {
	Key           string
	KafkaMetadata []kafka.TopicPartition
	Value         T
}

func NewMessage[T any](value T) *Message[T] {
	return &Message[T]{
		Value:         value,
		KafkaMetadata: []kafka.TopicPartition{},
	}
}

func NewEmptyMessage() *Message[any] {
	return &Message[any]{
		Value:         "",
		KafkaMetadata: []kafka.TopicPartition{},
	}
}

func (m *Message[T]) WithKafkaMetadata(tp kafka.TopicPartition) *Message[T] {
	m.KafkaMetadata = append(m.KafkaMetadata, tp)
	return m
}

func (m *Message[T]) WithKey(k string) *Message[T] {
	m.Key = k
	return m
}

func ReplaceValue[T, R any](m *Message[T], v R) *Message[R] {
	return &Message[R]{
		Key:           m.Key,
		KafkaMetadata: m.KafkaMetadata,
		Value:         v,
	}
}
