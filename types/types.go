package types

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

type KafkaMessage[T any] struct {
	Key            string
	TopicPartition kafka.TopicPartition
	Value          T
}