package operator

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/twoojoo/ttask/task"
	"github.com/twoojoo/ttask/types"
)

func fromKafka(consumer *kafka.Consumer, timeout ...time.Duration) task.Operator[any, types.KafkaMessage[[]byte]] {
	return func(m *task.Meta, x *task.Message[any], next *task.Step) {
		to := time.Second

		if len(timeout) > 0 {
			to = timeout[0]
		}

		for {
			msg, err := consumer.ReadMessage(to)
			if err == nil {

				m.ExecNext(task.NewMessage(types.KafkaMessage[[]byte]{
					TopicPartition: msg.TopicPartition,
					Key:            string(msg.Key),
					Value:          msg.Value,
				}), next)
			} else if !err.(kafka.Error).IsTimeout() {
				// TODO timeout error handling
				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		}
	}
}

// Source: trigger a Task execution for each received message.
func FromKafka(consumer *kafka.Consumer, timeout ...time.Duration) *task.TTask[any, types.KafkaMessage[[]byte]] {
	return task.T(task.Task[any](), fromKafka(consumer, timeout...))
}
