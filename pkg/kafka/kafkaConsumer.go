package kafka

import "github.com/Shopify/sarama"

func ListenFromKafka(topic string) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
}
