package operator

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/fatih/color"
	"github.com/twoojoo/ttask/task"
	"github.com/twoojoo/ttask/types"
)

func KafkaCommit[T any](consumer *kafka.Consumer) task.Operator[types.KafkaMessage[T], types.KafkaMessage[T]] {
	return func(m *task.Meta, x *task.Message[types.KafkaMessage[T]], next *task.Step) {
		_, err := consumer.CommitOffsets([]kafka.TopicPartition{x.Value.TopicPartition})
		if err != nil {
			log.Fatal(err)
		}

		m.ExecNext(x, next)
	}
}

func KafkaCommitMany[T any](consumer *kafka.Consumer) task.Operator[[]types.KafkaMessage[T], []types.KafkaMessage[T]] {
	return func(m *task.Meta, x *task.Message[[]types.KafkaMessage[T]], next *task.Step) {
		tp := []kafka.TopicPartition{}

		for _, v := range x.Value {
			tp = append(tp, v.TopicPartition)
		}

		_, err := consumer.CommitOffsets(tp)
		if err != nil {
			log.Fatal(err)
		}

		m.ExecNext(x, next)
	}
}

func PrintKafkaMessageMetadata[T any]() task.Operator[types.KafkaMessage[T], types.KafkaMessage[T]] {
	boldString := color.New(color.Bold).SprintFunc()
	header := boldString(color.GreenString("KafkaMessage"))
	decUp := boldString(color.GreenString("┏"))
	decMid := boldString(color.GreenString("┃"))
	decDown := boldString(color.GreenString("┗"))

	return func(m *task.Meta, x *task.Message[types.KafkaMessage[T]], next *task.Step) {
		log.Print(decUp, header)
		log.Printf("%s topic: \t  %s\n", decMid, *x.Value.TopicPartition.Topic)
		log.Printf("%s partition:  %v\n", decMid, x.Value.TopicPartition.Partition)
		log.Printf("%s offset: \t  %s\n", decMid, x.Value.TopicPartition.Offset)
		log.Printf("%s key: \t  %s\n", decDown, string(x.Value.Key))

		m.ExecNext(x, next)
	}
}

func PrintKafkaCommitMetadata[T any]() task.Operator[types.KafkaMessage[T], types.KafkaMessage[T]] {
	boldString := color.New(color.Bold).SprintFunc()
	header := boldString(color.CyanString("KafkaCommit"))
	decUp := boldString(color.CyanString("┏"))
	decMid := boldString(color.CyanString("┃"))
	decDown := boldString(color.CyanString("┗"))

	return func(m *task.Meta, x *task.Message[types.KafkaMessage[T]], next *task.Step) {
		// maxOffset := x.Value.TopicPartition.Offset
		// topicsSet := map[string]struct{}{}

		// for _, md := range x.Value.TopicPartition {
		// 	if md.Offset > maxOffset {
		// 		maxOffset = md.Offset
		// 	}

		// 	topicsSet[*md.Topic] = struct{}{}
		// }

		// topicsStr := ""

		// first := true
		// for k := range topicsSet {
		// 	if first {
		// 		topicsStr += k
		// 		first = false
		// 	} else {
		// 		topicsStr += ", "
		// 		topicsStr += k
		// 	}
		// }
		log.Print(decUp, header)
		log.Printf("%s topic: \t  %s\n", decMid, *x.Value.TopicPartition.Topic)
		log.Printf("%s partition:  %v\n", decMid, x.Value.TopicPartition.Partition)
		log.Printf("%s offset: \t  %s\n", decMid, x.Value.TopicPartition.Offset)
		log.Printf("%s key: \t  %s\n", decDown, string(x.Value.Key))

		m.ExecNext(x, next)
	}
}
