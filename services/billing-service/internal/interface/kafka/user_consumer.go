package kafka

import (
	"context"
	"log/slog"
	"time"

	"github.com/Romasmi/e-learning-arhitecture/billing-service/internal/domain/account"
	"github.com/Romasmi/e-learning-arhitecture/billing-service/internal/infrastructure/db/postgres"
	userapi "github.com/Romasmi/e-learning-arhitecture/gen/go/user"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/encoding/protojson"
)

type UserConsumer struct {
	reader *kafka.Reader
	repo   *postgres.AccountRepository
}

func NewUserConsumer(brokers []string, topic string, repo *postgres.AccountRepository) *UserConsumer {
	return &UserConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  "billing-service",
			MinBytes: 1,
			MaxBytes: 10e6, // 10MB
			MaxWait:  500 * time.Millisecond,
		}),
		repo: repo,
	}
}

func (c *UserConsumer) Start(ctx context.Context) {
	slog.Info("Starting user.created consumer")
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			slog.Error("failed to read message", "error", err)
			continue
		}

		var event userapi.UserCreated
		if err := protojson.Unmarshal(m.Value, &event); err != nil {
			slog.Error("failed to unmarshal user.created event", "error", err)
			continue
		}

		userID, err := uuid.Parse(event.User.Id)
		if err != nil {
			slog.Error("failed to parse user_id from event", "error", err, "user_id", event.User.Id)
			continue
		}

		acc := &account.Account{
			UserID:  userID,
			Balance: 0,
		}

		if err := c.repo.CreateAccount(ctx, acc); err != nil {
			slog.Error("failed to create account for user", "error", err, "user_id", userID)
			continue
		}

		slog.Info("Account created for user", "user_id", userID)
	}
}

func (c *UserConsumer) Close() error {
	return c.reader.Close()
}
