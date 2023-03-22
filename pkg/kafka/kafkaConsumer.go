package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/gommon/log"
)

func ListenFromKafka(topic string) []byte {
	fmt.Printf("Starting consumer...")
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Errorf("Something went wrong: %v", err)
	}

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Errorf("Something went wrong: %v", err)
	}

	msg, err := c.ReadMessage(-1)
	if err == nil {
		data := string(msg.Value)
		log.Printf("Message on %s: %s\n", msg.TopicPartition, data)
	} else {
		// The client will automatically try to recover from all errors.
		log.Errorf("Consumer error: %v (%v)\n", err, msg)
	}

	/*	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			data := string(msg.Value)
			log.Printf("Message on %s: %s\n", msg.TopicPartition, data)
		} else {
			// The client will automatically try to recover from all errors.
			log.Errorf("Consumer error: %v (%v)\n", err, msg)
		}
	}*/

	err = c.Close()
	if err != nil {
		log.Errorf("Something went wrong: %v", err)
	}

	return msg.Value
}
