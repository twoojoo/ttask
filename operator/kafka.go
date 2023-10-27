package operator

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/twoojoo/ttask/task"
)

func KafkaCommit[T any](consumer *kafka.Consumer) task.Operator[T, T] {
	return func(m *task.Meta, x *task.Message[T], next *task.Step) {

		_, err := consumer.CommitOffsets(x.KafkaMetadata)
		if err != nil {
			log.Fatal(err)
		}

		m.ExecNext(x, next)
	}
}
