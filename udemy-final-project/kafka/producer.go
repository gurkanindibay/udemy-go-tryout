package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

// Producer handles publishing messages to Kafka
type Producer struct {
	writer *kafka.Writer
	topic  string
}

// EventMessage represents the structure of messages sent to Kafka
type EventMessage struct {
	Action string      `json:"action"`
	Event  interface{} `json:"event"`
}

// NewProducer creates a new Kafka producer instance
func NewProducer() (*Producer, error) {
	brokers := []string{getEnv("KAFKA_BROKERS", "localhost:9092")}
	topic := "events"

	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{
		writer: writer,
		topic:  topic,
	}, nil
}

// PublishEvent sends an event message to Kafka with the specified action and event ID
func (p *Producer) PublishEvent(action string, eventID string, event interface{}) error {
	message := EventMessage{
		Action: action,
		Event:  event,
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	key := action + "-" + eventID // Combine action and eventID for partitioning

	err = p.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: jsonMessage,
		},
	)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		return err
	}

	log.Printf("Published event to Kafka: %s for event ID %s", action, eventID)
	return nil
}

// Close closes the Kafka producer and releases resources
func (p *Producer) Close() error {
	return p.writer.Close()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
