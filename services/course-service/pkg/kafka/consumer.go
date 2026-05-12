package kafka

import (
	"context"
	"fmt"

	kgo "github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kgo.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	return &Consumer{
		reader: kgo.NewReader(kgo.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
	}
}

func (c *Consumer) Consume(ctx context.Context, handler func(msg kgo.Message) error) {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			fmt.Printf("could not read message: %v\n", err)
			break
		}
		if err := handler(m); err != nil {
			fmt.Printf("error handling message: %v\n", err)
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
