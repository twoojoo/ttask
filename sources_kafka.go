package ttask

import (
	"errors"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/fatih/color"
)

func fromKafka(consumer *kafka.Consumer, logger bool, timeout ...time.Duration) Operator[any, KafkaMessage[[]byte]] {
	return func(inner *Inner, x *Message[any], next *Step) {
		to := time.Second

		if len(timeout) > 0 {
			to = timeout[0]
		}

		for {
			msg, err := consumer.ReadMessage(to)
			if err == nil {
				tMsg := newMessage(KafkaMessage[[]byte]{
					TopicPartition: msg.TopicPartition,
					Key:            string(msg.Key),
					Value:          msg.Value,
				}).withKey(string(msg.Key))

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
func FromKafka(taskId string, consumer *kafka.Consumer, logger bool, timeout ...time.Duration) *TTask[any, KafkaMessage[[]byte]] {
	return Via(Task[any](taskId), fromKafka(consumer, logger, timeout...))
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
