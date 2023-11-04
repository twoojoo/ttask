package ttask

import (
	"errors"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/fatih/color"
)

type KafkaSinkOpts struct {
	ContinueOnError    bool
	SkipErrorReporting bool
	Logger             bool
}

// Sink: send the message to a kafka topic
func ToKafka[T any](producer *kafka.Producer, topic string, toBytes func(x T) []byte, options KafkaSinkOpts) Operator[T, T] {
	return func(inner *Inner, x *Message[T], next *Step) {
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
			inner.Error(errors.New("TTask [ToKafka] error: " + err.Error()))
			return
		}

		event := <-ch

		switch ev := event.(type) {
		case *kafka.Message:
			inner.ExecNext(x, next)

			if options.Logger {
				logKafkaSend(x.Key, ev.TopicPartition)
			}
		default:
			if options.ContinueOnError {
				inner.ExecNext(x, next)
			}

			if !options.SkipErrorReporting {
				inner.Error(errors.New("TTask [ToKafka] error: " + ev.String()))
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
