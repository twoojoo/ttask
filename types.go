package ttask

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"golang.org/x/exp/constraints"
)

type KafkaMessage[T any] struct {
	Key            string
	TopicPartition kafka.TopicPartition
	Value          T
}

type Number interface {
	constraints.Integer | constraints.Float
}
