package sink

import (
	"errors"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/twoojoo/ttask/task"
)

type KafkaSinkOpts struct {
	ContinueOnError    bool
	SkipErrorReporting bool
}

// Sink: send the message to a kafka topic
func ToKafka[T any](producer *kafka.Producer, topic string, toBytes func(x T) []byte, options KafkaSinkOpts) task.Operator[T, T] {
	return func(m *task.Meta, x *task.Message[T], next *task.Step) {
		ch := make(chan kafka.Event)

		producer.Produce(&kafka.Message{
			Key:   []byte(x.Key),
			Value: toBytes(x.Value),
		}, ch)

		event := <-ch

		switch ev := event.(type) {
		case *kafka.Message:
			m.ExecNext(x, next)
		default:
			if options.ContinueOnError {
				m.ExecNext(x, next)
			}

			if !options.SkipErrorReporting {
				m.Error(errors.New(ev.String()))
			}
		}
	}
}
