package main

import (
	"context"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	. "github.com/twoojoo/ttask"
)

func main() {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  "localhost",
		"group.id":           "test-kafka-0",
		"enable.auto.commit": "false",
	})

	if err != nil {
		log.Fatal(err)
	}

	consumer.Subscribe("sp-gpcs-reservations-raw", func(c *kafka.Consumer, e kafka.Event) error {
		log.Println(e.String())
		return nil
	})

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
	})

	if err != nil {
		log.Fatal(err)
	}

	logs := true

	T(T(T(T(
		FromKafka("t1", consumer, logs, time.Minute),
		WithEventTime(func(_ KafkaMessage[[]byte]) time.Time {
			log.Println("#> extracting event time..")
			return time.Now().Add(-time.Second)
		})),
		WithCustomKey(func(_ KafkaMessage[[]byte]) string {
			log.Println("#> setting custom key..")
			return "c-key"
		})),
		KafkaCommit[[]byte](consumer, logs)),
		ToKafka[KafkaMessage[[]byte]](producer,
			"sp-gpcs-reservations-parsed",
			func(x KafkaMessage[[]byte]) []byte {
				return x.Value
			},
			KafkaSinkOpts{Logger: logs},
		),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Run(context.Background())
}
