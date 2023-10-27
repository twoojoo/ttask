package main

import (
	"context"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	. "github.com/twoojoo/ttask/operator"
	. "github.com/twoojoo/ttask/source"
	. "github.com/twoojoo/ttask/task"
)

func main() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  "localhost",
		"group.id":           "test-kafka",
		"enable.auto.commit": "false",
	})

	if err != nil {
		log.Fatal(err)
	}

	c.Subscribe("sp-gpcs-reservations-raw", func(c *kafka.Consumer, e kafka.Event) error {
		log.Println(e.String())
		return nil
	})

	T(T(FromKafka(c),
		Print[*kafka.Message]("received >")),
		Tap(func(x *kafka.Message) {
			fmt.Println("-------------------------")
			fmt.Println("key:", string(x.Key))
			fmt.Println("topic:", *x.TopicPartition.Topic)
			fmt.Println("offset:", x.TopicPartition.Offset)
			fmt.Println("-------------------------")
		})).
		Catch(func(m *Meta, e error) {
			v := m.Context.Value("k1").(string)
			log.Println("ctx value was:", v)
			log.Println(e)
		}).Run(context.Background())
}
