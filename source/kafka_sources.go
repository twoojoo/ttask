package operator

import (
	"errors"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/fatih/color"
	"github.com/twoojoo/ttask/task"
	"github.com/twoojoo/ttask/types"
)

func fromKafka(consumer *kafka.Consumer, logger bool, timeout ...time.Duration) task.Operator[any, types.KafkaMessage[[]byte]] {
	return func(inner *task.Inner, x *task.Message[any], next *task.Step) {
		to := time.Second

		if len(timeout) > 0 {
			to = timeout[0]
		}

		for {
			msg, err := consumer.ReadMessage(to)
			if err == nil {
				tMsg := task.NewMessage(types.KafkaMessage[[]byte]{
					TopicPartition: msg.TopicPartition,
					Key:            string(msg.Key),
					Value:          msg.Value,
				}).WithKey(string(msg.Key))

				if logger {
					logKafkaMessage(msg)
				}

				inner.ExecNext(tMsg, next)
			} else {
				inner.Error(errors.New("TTask [FromKafka] error: " + err.Error()))
			}
		}
	}
}

// Source: trigger a Task execution for each received message.
func FromKafka(taskId string, consumer *kafka.Consumer, logger bool, timeout ...time.Duration) *task.TTask[any, types.KafkaMessage[[]byte]] {
	return task.T(task.Task[any](taskId), fromKafka(consumer, logger, timeout...))
}

func logKafkaMessage(msg *kafka.Message) {
	boldString := color.New(color.Bold).SprintFunc()
	header := boldString(color.GreenString("KafkaMessageReceived"))
	decUp := boldString(color.GreenString("┏"))
	decMid := boldString(color.GreenString("┃"))
	decDown := boldString(color.GreenString("┗"))

	log.Print(decUp, header)
	log.Printf("%s topic: \t  %s\n", decMid, *msg.TopicPartition.Topic)
	log.Printf("%s partition:  %v\n", decMid, msg.TopicPartition.Partition)
	log.Printf("%s offset: \t  %s\n", decMid, msg.TopicPartition.Offset)
	log.Printf("%s key: \t  %s\n", decDown, string(msg.Key))
}
