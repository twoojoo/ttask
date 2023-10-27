package operator

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/twoojoo/ttask/task"
)

func fromKafka(consumer *kafka.Consumer, timeout ...time.Duration) task.Operator[any, []byte] {
	return func(m *task.Meta, x *task.Message[any], next *task.Step) {
		to := time.Second

		if len(timeout) > 0 {
			to = timeout[0]
		}

		for {
			msg, err := consumer.ReadMessage(to)
			if err == nil {
				taskMsg := task.NewMessage(msg.Value).
					WithKafkaMetadata(msg.TopicPartition).
					WithKey(string(msg.Key))

				m.ExecNext(taskMsg, next)
			} else if !err.(kafka.Error).IsTimeout() {
				// TODO timeout error handling
				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		}
	}
}


//Source: trigger a Task execution for each received message.
func FromKafka(consumer *kafka.Consumer, timeout ...time.Duration) *task.TTask[any, []byte] {
	return task.T(task.Task[any](), fromKafka(consumer, timeout...))
}
