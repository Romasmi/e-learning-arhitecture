package postgres

import (
	"context"
	"fmt"

	"github.com/Romasmi/e-learning-arhitecture/notification-service/internal/domain/message"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepository struct {
	pool *pgxpool.Pool
}

func NewMessageRepository(pool *pgxpool.Pool) *MessageRepository {
	return &MessageRepository{pool: pool}
}

func (r *MessageRepository) CreateMessage(ctx context.Context, m *message.Message) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO messages (id, user_id, order_id, type, created_at, updated_at) VALUES ($1, $2, $3, $4, NOW(), NOW())", m.ID, m.UserID, m.OrderID, m.Type)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}
	return nil
}

func (r *MessageRepository) ListMessages(ctx context.Context, userID uuid.UUID) ([]*message.Message, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, user_id, order_id, type, created_at, updated_at FROM messages WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}
	defer rows.Close()

	var messages []*message.Message
	for rows.Next() {
		m := &message.Message{}
		if err := rows.Scan(&m.ID, &m.UserID, &m.OrderID, &m.Type, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, m)
	}
	return messages, nil
}
