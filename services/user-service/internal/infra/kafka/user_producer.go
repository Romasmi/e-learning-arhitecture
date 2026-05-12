package kafka

import (
	"context"
	"fmt"
	"time"

	api "github.com/Romasmi/e-learning-arhitecture/gen/go/user"
	"github.com/Romasmi/e-learning-arhitecture/user-service/internal/domain/user"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserProducer struct {
	Writer *kafka.Writer
}

func NewUserProducer(brokers []string, topic string) *UserProducer {
	return &UserProducer{
		Writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond,
			RequiredAcks: kafka.RequireAll,
		},
	}
}

func (p *UserProducer) EmitUserCreated(ctx context.Context, u *user.User) error {
	event := &api.UserCreated{
		EventId:    uuid.New().String(),
		OccurredAt: timestamppb.New(time.Now()),
		User: &api.User{
			Id:        u.ID.String(),
			Username:  u.Username,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Phone:     u.Phone,
			Password:  u.Password,
			CreatedAt: timestamppb.New(u.CreatedAt),
		},
	}

	payload, err := protojson.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = p.Writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(u.ID.String()),
		Value: payload,
	})
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

func (p *UserProducer) Close() error {
	return p.Writer.Close()
}
