package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

// SendToKafka take a topic name and message with format of []byte
func SendToKafka(topic string, message []byte) {
	// to produce messages
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatalf("failed to dial leader: %v", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: message},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
