package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
	topic  string
}

type EventMessage struct {
	Action string      `json:"action"`
	Event  interface{} `json:"event"`
}

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

func (p *Producer) PublishEvent(action string, event interface{}) error {
	message := EventMessage{
		Action: action,
		Event:  event,
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = p.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(action),
			Value: jsonMessage,
		},
	)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		return err
	}

	log.Printf("Published event to Kafka: %s", action)
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}