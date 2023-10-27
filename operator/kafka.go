package operator

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/fatih/color"
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

func PrintKafkaMessageMetadata() task.Operator[[]byte, []byte] {
	boldString := color.New(color.Bold).SprintFunc()
	header := boldString(color.GreenString("KafkaMessage"))
	decUp := boldString(color.GreenString("┏"))
	decMid := boldString(color.GreenString("┃"))
	decDown := boldString(color.GreenString("┗"))

	return func(m *task.Meta, x *task.Message[[]byte], next *task.Step) {
		log.Print(decUp, header)
		log.Printf("%s topic: \t  %s\n", decMid, *x.KafkaMetadata[0].Topic)
		log.Printf("%s partition:  %v\n", decMid, *&x.KafkaMetadata[0].Partition)
		log.Printf("%s offset: \t  %s\n", decMid, x.KafkaMetadata[0].Offset)
		log.Printf("%s key: \t  %s\n", decDown, string(x.Key))

		m.ExecNext(x, next)
	}
}

func PrintKafkaCommitMetadata[T any]() task.Operator[T, T] {
	boldString := color.New(color.Bold).SprintFunc()
	header := boldString(color.CyanString("KafkaCommit"))
	decUp := boldString(color.CyanString("┏"))
	decMid := boldString(color.CyanString("┃"))
	decDown := boldString(color.CyanString("┗"))

	return func(m *task.Meta, x *task.Message[T], next *task.Step) {
		maxOffset := x.KafkaMetadata[0].Offset
		topicsSet := map[string]struct{}{}

		for _, md := range x.KafkaMetadata {
			if md.Offset > maxOffset {
				maxOffset = md.Offset
			}

			topicsSet[*md.Topic] = struct{}{}
		}

		topicsStr := ""

		first := true
		for k := range topicsSet {
			if first {
				topicsStr += k
				first = false
			} else {
				topicsStr += ", "
				topicsStr += k
			}
		}

		log.Print(decUp, header)
		log.Printf("%s topics: \t  %s\n", decMid, topicsStr)
		log.Printf("%s partition:  %v\n", decMid, *&x.KafkaMetadata[0].Partition)
		log.Printf("%s max offset: %s\n", decMid, maxOffset)
		log.Printf("%s key: \t  %s\n", decDown, string(x.Key))

		m.ExecNext(x, next)
	}
}
