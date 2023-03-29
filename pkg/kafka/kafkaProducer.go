package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/gommon/log"
)

type ProducerKafka struct {
	Producer *kafka.Producer
	Topic    string
}

func NewProducerKafka(producer *kafka.Producer, topic string) *ProducerKafka {
	return &ProducerKafka{
		Producer: producer,
		Topic:    topic,
	}
}

// SendToKafka take a topic name and message with format of []byte
/*func SendToKafka(topic string, message []byte) error {

	// Producer
	log.Print("Starting producer...")
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Errorf("Cannot create a producer: %v", err)
		return err
	}
	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Errorf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					log.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)
	if err != nil {
		log.Errorf("Something went wrong: %v", err)
		return err
	}

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)

	return nil
}*/

func (p *ProducerKafka) SendToKafkaWithMessage(message []byte) error {
	// Delivery report handler for produced messages
	go func() {
		for e := range p.Producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Errorf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					log.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic
	err := p.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.Topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)
	if err != nil {
		log.Errorf("Something went wrong: %v", err)
		return err
	}

	// Wait for message deliveries before shutting down
	p.Producer.Flush(15 * 1000)

	return nil
}
