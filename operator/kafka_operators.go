package operator

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/fatih/color"
	"github.com/twoojoo/ttask/task"
	"github.com/twoojoo/ttask/types"
)

//Perform a commit on the current kafka message.
func KafkaCommit[T any](consumer *kafka.Consumer, logger bool) task.Operator[types.KafkaMessage[T], types.KafkaMessage[T]] {
	return func(m *task.Inner, x *task.Message[types.KafkaMessage[T]], next *task.Step) {
		tp := x.Value.TopicPartition

		_, err := consumer.CommitOffsets([]kafka.TopicPartition{tp})
		if err != nil {
			m.Error(err)
			return
		}

		if logger {
			logKafkaCommit(tp)
		}

		m.ExecNext(x, next)
	}
}

// func KafkaCommitMany[T any](consumer *kafka.Consumer, logger bool) task.Operator[[]types.KafkaMessage[T], []types.KafkaMessage[T]] {
// 	return func(m *task.Inner, x *task.Message[[]types.KafkaMessage[T]], next *task.Step) {
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
