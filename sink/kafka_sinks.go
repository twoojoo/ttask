package sink

import (
	"errors"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/fatih/color"
	"github.com/twoojoo/ttask/task"
)

type KafkaSinkOpts struct {
	ContinueOnError    bool
	SkipErrorReporting bool
	Logger             bool
}

// Sink: send the message to a kafka topic
func ToKafka[T any](producer *kafka.Producer, topic string, toBytes func(x T) []byte, options KafkaSinkOpts) task.Operator[T, T] {
	return func(m *task.Inner, x *task.Message[T], next *task.Step) {
		ch := make(chan kafka.Event)

		err := producer.Produce(&kafka.Message{
			Key:   []byte(x.Key),
			Value: toBytes(x.Value),
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
		}, ch)

		if err != nil {
			m.Error(errors.New("TTask [ToKafka] error: " + err.Error()))
			return
		}

		event := <-ch

		switch ev := event.(type) {
		case *kafka.Message:
			m.ExecNext(x, next)

			if options.Logger {
				logKafkaSend(x.Key, ev.TopicPartition)
			}
		default:
			if options.ContinueOnError {
				m.ExecNext(x, next)
			}

			if !options.SkipErrorReporting {
				m.Error(errors.New("TTask [ToKafka] error: " + ev.String()))
				return
			}
		}
	}
}

func logKafkaSend(key string, tp kafka.TopicPartition) {
	boldString := color.New(color.Bold).SprintFunc()
	header := boldString(color.YellowString("KafkaMessageProduced"))
	decUp := boldString(color.YellowString("┏"))
	decMid := boldString(color.YellowString("┃"))
	decDown := boldString(color.YellowString("┗"))

	log.Print(decUp, header)
	log.Printf("%s topic: \t  %s\n", decMid, *tp.Topic)
	log.Printf("%s partition:  %v\n", decMid, tp.Partition)
	log.Printf("%s offset: \t  %s\n", decMid, tp.Offset)
	log.Printf("%s key: \t  %s\n", decDown, key)
}
