package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/gommon/log"
	"time"
)

type ProducerKafka struct {
	Producer *kafka.Producer
}

func NewProducerKafka(kafkaHost string) *ProducerKafka {
	// To create kafka producer as a 'ProducerKafka' struct
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaHost})
	if err != nil {
		log.Errorf("Cannot create a producer: %v", err)
	}

	return &ProducerKafka{
		Producer: p,
	}
}

func (p *ProducerKafka) SendToKafkaWithMessage(message []byte, topic string) error {
	// Delivery report handler for produced messages
	go func() {
		for e := range p.Producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Errorf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					log.Infof("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	/*
		msg := &kafka.Message{
		    TopicPartition: kafka.TopicPartition{
		        Topic:     &topic,
		        Partition: kafka.PartitionAny,
		    },
		    Key:   []byte("123"), // Partition key value
		    Value: message,
		}
		Burada Key alanına "123" değeri atanmıştır. Bu, mesajın "123" değerine sahip partition key'e sahip olan bir partition'a yazılmasını sağlayacaktır. Partition key değeri, mesajın içeriğine göre değişebilir veya sabit bir değer olarak belirlenebilir.
	*/

	// Produce messages to topic
	err := p.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)
	if err != nil {
		log.Errorf("Something went wrong: %v", err)
		return err
	}

	// Wait for message deliveries before shutting down
	p.Producer.Flush(15 * (int(time.Millisecond)))

	return nil
}
