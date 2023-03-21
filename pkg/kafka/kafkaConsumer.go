package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/gommon/log"
)

func ListenFromKafka(topic string) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "my-group",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Error(err)
	}

	defer c.Close()

	// consumer subscribe
	err = c.SubscribeTopics([]string{topic, "^aRegex.*[Tt]opic"}, nil)
	if err != nil {
		log.Error(err)
	}

	// to read message
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Received message: %s\n", string(msg.Value))
		} else {
			// handle error
			fmt.Printf("Error while consuming message: %v\n", err)
		}
	}
}
