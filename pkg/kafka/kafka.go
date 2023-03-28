package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
	"time"
)

type KafkaConsumer struct {
	consumer    *kafka.Consumer
	lastMessage kafka.Message
	logger      *logrus.Logger
}

func NewKafkaConsumer(prd *kafka.Consumer, logger *logrus.Logger) *KafkaConsumer {
	return &KafkaConsumer{
		consumer:    prd,
		lastMessage: kafka.Message{},
		logger:      logger,
	}
}

func (c *KafkaConsumer) SubscribeToTopics(topics []string) error {
	err := c.consumer.SubscribeTopics(topics, nil)
	return err
}

func (c *KafkaConsumer) ConsumeFromTopics(bulkConsumeIntervalInSeconds int64, bulkConsumeMaxTimeoutInSeconds int, maxReadCount int) ([]kafka.Message, error) {
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
				c.logger.Errorf("Kafka read messages failed. | Error: %v\n", err)
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

func (c *KafkaConsumer) AckLastMessage() {
	if &c.lastMessage != nil {
		_, err := c.consumer.CommitMessage(&c.lastMessage)
		if err != nil {
			c.logger.Errorf("Ack Last message failed. | Error: %v\n", err)
		}
	}
}
