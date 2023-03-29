package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/gommon/log"
	"time"
)

type ConsumerKafka struct {
	consumer    *kafka.Consumer
	lastMessage kafka.Message
}

func NewConsumerKafka(consumer *kafka.Consumer) *ConsumerKafka {
	return &ConsumerKafka{
		consumer:    consumer,
		lastMessage: kafka.Message{},
	}
}

func (c *ConsumerKafka) SubscribeToTopics(topics []string) error {
	err := c.consumer.SubscribeTopics(topics, nil)
	return err
}

func (c *ConsumerKafka) ConsumeFromTopics(bulkConsumeIntervalInSeconds int64, bulkConsumeMaxTimeoutInSeconds int, maxReadCount int) ([]kafka.Message, error) {
	messages := make([]kafka.Message, 0)
	timeoutCount := 0
	start := time.Now()
	for {
		msg, err := c.consumer.ReadMessage(time.Duration(bulkConsumeMaxTimeoutInSeconds) * time.Second)

		elapsedTime := time.Since(start)
		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrTimedOut {
				timeoutCount++
				if timeoutCount > 10 {
					return messages, nil
				}
				continue
			} else {
				log.Errorf("Kafka read messages failed. | Error: %v\n", err)
				if len(messages) > 0 {
					return messages, err
				} else {
					continue
				}
			}
		}

		if msg != nil {
			c.lastMessage = *msg
			messages = append(messages, *msg)
		}

		if elapsedTime.Milliseconds() > (bulkConsumeIntervalInSeconds*1000) || len(messages) >= maxReadCount {
			return messages, nil
		}
	}
}

func (c *ConsumerKafka) AckLastMessage() {
	if &c.lastMessage != nil {
		_, err := c.consumer.CommitMessage(&c.lastMessage)
		if err != nil {
			log.Errorf("Ack Last message failed. | Error: %v\n", err)
		}
	}
}

func (c *ConsumerKafka) ListenFromKafkaWithoutTopic(topic string) ([]byte, error) {
	err := c.consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Errorf("Something went wrong: %v", err)
	}
	defer c.consumer.Close()

	for {
		msg, err := c.consumer.ReadMessage(-1)
		if err == nil {
			data := string(msg.Value)
			log.Printf("Message on %s: %s\n", msg.TopicPartition, data)
			return msg.Value, nil
		} else {
			// The client will automatically try to recover from all errors.
			log.Errorf("Consumer error: %v (%v)\n", err, msg)
			return []byte{}, err
		}
	}
}

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

	err = c.Close()
	if err != nil {
		log.Errorf("Something went wrong: %v", err)
	}

	return msg.Value
}
