package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	topic  string
}

func NewConsumer() (*Consumer, error) {
	brokers := []string{getEnv("KAFKA_BROKERS", "localhost:9092")}
	topic := "events"

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  "event-consumer-group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &Consumer{
		reader: reader,
		topic:  topic,
	}, nil
}

func (c *Consumer) StartConsuming() {
	log.Println("Kafka consumer started, waiting for messages...")

	for {
		m, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Consumer error: %v", err)
			continue
		}

		c.processMessage(&m)
	}
}

func (c *Consumer) processMessage(msg *kafka.Message) {
	var eventMessage EventMessage
	err := json.Unmarshal(msg.Value, &eventMessage)
	if err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	log.Printf("Received event: Action=%s, Event=%+v", eventMessage.Action, eventMessage.Event)

	// Here you can add logic to handle the event, e.g., send notifications, update analytics, etc.
	// For this demo, we'll just log it.
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}