package main

import (
	"context"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	. "github.com/twoojoo/ttask"
)

func main() {
	c1, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  "localhost",
		"group.id":           "test-kafka-1",
		"enable.auto.commit": "false",
	})

	if err != nil {
		log.Fatal(err)
	}

	c1.Subscribe("sp-gpcs-reservations-parsed", func(c *kafka.Consumer, e kafka.Event) error {
		log.Println(e.String())
		return nil
	})

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
	})

	if err != nil {
		log.Fatal(err)
	}

	T(
		FromKafka("t1", c1, true, time.Minute),
		KafkaCommit[[]byte](c1, true),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Run(context.Background())

	topic := "sp-gpcs-reservations-raw"
	key := "k1"
	p.Produce(&kafka.Message{
		Key:            []byte(key),
		TopicPartition: kafka.TopicPartition{Topic: &topic},
		Value:          []byte("this is a message"),
	}, nil)
}
