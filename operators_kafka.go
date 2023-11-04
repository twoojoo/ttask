package ttask

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/fatih/color"
)

// Perform a commit on the current kafka message.
func KafkaCommit[T any](consumer *kafka.Consumer, logger bool) Operator[KafkaMessage[T], KafkaMessage[T]] {
	return func(inner *Inner, x *Message[KafkaMessage[T]], next *Step) {
		tp := x.Value.TopicPartition

		_, err := consumer.CommitOffsets([]kafka.TopicPartition{tp})
		if err != nil {
			inner.Error(err)
			return
		}

		if logger {
			logKafkaCommit(tp)
		}

		inner.ExecNext(x, next)
	}
}

// func KafkaCommitMany[T any](consumer *kafka.Consumer, logger bool) Operator[[]types.KafkaMessage[T], []types.KafkaMessage[T]] {
// 	return func(inner *Inner, x *Message[[]types.KafkaMessage[T]], next *Step) {
// 		tp := []kafka.TopicPartition{}

// 		for _, v := range x.Value {
// 			tp = append(tp, v.TopicPartition)
// 		}

// 		_, err := consumer.CommitOffsets(tp)
// 		if err != nil {
// 			m.Error(err)
// 			return
// 		}

// 		if logger {

// 		}

// 		m.ExecNext(x, next)
// 	}
// }

func logKafkaCommit(tp kafka.TopicPartition) {
	boldString := color.New(color.Bold).SprintFunc()
	header := boldString(color.CyanString("KafkaOffsetCommited"))
	decUp := boldString(color.CyanString("┏"))
	decMid := boldString(color.CyanString("┃"))
	decDown := boldString(color.CyanString("┗"))

	log.Print(decUp, header)
	log.Printf("%s topic: \t  %s\n", decMid, *tp.Topic)
	log.Printf("%s partition:  %v\n", decMid, tp.Partition)
	log.Printf("%s offset: \t  %s\n", decDown, tp.Offset)
}
