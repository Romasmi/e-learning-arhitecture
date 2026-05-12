package kafka

import (
	"context"
	"log/slog"
	"time"

	//orderapi "github.com/Romasmi/e-learning-arhitecture/gen/go/order"
	// "github.com/Romasmi/e-learning-arhitecture/notification-service/internal/domain/message"
	"github.com/Romasmi/e-learning-arhitecture/notification-service/internal/infrastructure/db/postgres"
	"github.com/segmentio/kafka-go"
	//	"google.golang.org/protobuf/encoding/protojson"
)

type OrderConsumer struct {
	reader *kafka.Reader
	repo   *postgres.MessageRepository
}

func NewOrderConsumer(brokers []string, topic string, repo *postgres.MessageRepository) *OrderConsumer {
	return &OrderConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  "notification-service",
			MinBytes: 1,
			MaxBytes: 10e6,
			MaxWait:  500 * time.Millisecond,
		}),
		repo: repo,
	}
}

func (c *OrderConsumer) Start(ctx context.Context) {
	slog.Info("Starting order.placed consumer")
	for {
		_, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			slog.Error("failed to read message", "error", err)
			continue
		}

		//var event orderapi.OrderPlaced
		//if err := protojson.Unmarshal(m.Value, &event); err != nil {
		//	slog.Error("failed to unmarshal order.placed event", "error", err)
		//	continue
		//}
		//
		//msgType := "ORDER_FAILED"
		//if event.PaymentResult.Success {
		//	msgType = "ORDER_SUCCESS"
		//}
		//
		//msg := &message.Message{
		//	ID:        event.EventId,
		//	UserID:    event.User.Id,
		//	OrderID:   event.Order.Id,
		//	Type:      msgType,
		//	CreatedAt: time.Now(),
		//	UpdatedAt: time.Now(),
		//}
		//
		//if err := c.repo.CreateMessage(ctx, msg); err != nil {
		//	slog.Error("failed to persist message", "error", err, "order_id", event.Order.Id)
		//	continue
		//}
		//
		//slog.Info("Notification processed",
		//	"userId", event.User.Id,
		//	"orderId", event.Order.Id,
		//	"type", msgType,
		//	"timestamp", msg.CreatedAt,
		//)
	}
}

func (c *OrderConsumer) Close() error {
	return c.reader.Close()
}
