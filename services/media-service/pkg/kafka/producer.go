package kafka

import (
	"context"
	"encoding/json"
	"time"

	kgo "github.com/segmentio/kafka-go"
)

type Event struct {
	EventType  string          `json:"event_type"`
	AssetID    string          `json:"asset_id"`
	LessonID   string          `json:"lesson_id"`
	Payload    json.RawMessage `json:"payload"`
	OccurredAt time.Time       `json:"occurred_at"`
}

type Producer struct {
	writer *kgo.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kgo.Writer{
			Addr:     kgo.TCP(brokers...),
			Topic:    topic,
			Balancer: &kgo.LeastBytes{},
		},
	}
}

func (p *Producer) PublishAsync(event Event) {
	go func() {
		payload, err := json.Marshal(event)
		if err != nil {
			return
		}

		_ = p.writer.WriteMessages(context.Background(), kgo.Message{
			Key:   []byte(event.AssetID),
			Value: payload,
		})
	}()
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
