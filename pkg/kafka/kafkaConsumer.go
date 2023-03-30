package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/gommon/log"
	"time"
)

type ConsumerKafka struct {
	Consumer    *kafka.Consumer
	LastMessage kafka.Message
}

func NewConsumerKafka() *ConsumerKafka {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Errorf("Kafka consumer didn't work. Error:%v", err)
	}
	return &ConsumerKafka{
		Consumer:    c,
		LastMessage: kafka.Message{},
	}
}

func (c *ConsumerKafka) SubscribeToTopics(topics []string) error {
	err := c.Consumer.SubscribeTopics(topics, nil)
	return err
}

// bulkConsumeIntervalInSeconds: bulk reading interval (in seconds)
// bulkConsumeMaxTimeoutInSeconds: maximum read time (in seconds)
// maxReadCount: maximum number of messages to read

func (c *ConsumerKafka) ConsumeFromTopics(bulkConsumeIntervalInSeconds int64, bulkConsumeMaxTimeoutInSeconds int, maxReadCount int) ([]kafka.Message, error) {

	messages := make([]kafka.Message, 0)
	timeoutCount := 0
	start := time.Now()

	for {
		msg, err := c.Consumer.ReadMessage(time.Duration(bulkConsumeMaxTimeoutInSeconds) * time.Second)

		elapsedTime := time.Since(start)
		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrTimedOut {
				timeoutCount++
				if timeoutCount > 2 {
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
			c.LastMessage = *msg
			messages = append(messages, *msg)
		}

		if elapsedTime.Milliseconds() > (bulkConsumeIntervalInSeconds*1000) || len(messages) >= maxReadCount {
			return messages, nil
		}
	}
}

func (c *ConsumerKafka) AckLastMessage() {
	if &c.LastMessage != nil {
		_, err := c.Consumer.CommitMessage(&c.LastMessage)
		if err != nil {
			log.Errorf("Ack Last message failed. | Error: %v\n", err)
		}
	}
}

/*func ListenFromKafka(topic string) []byte {
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
}*/
