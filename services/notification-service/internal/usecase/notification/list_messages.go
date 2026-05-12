package notification

import (
	"context"

	"github.com/Romasmi/e-learning-arhitecture/notification-service/internal/domain/message"
	"github.com/google/uuid"
)

type ListMessagesUseCase struct {
	repo Repository
}

func NewListMessagesUseCase(repo Repository) *ListMessagesUseCase {
	return &ListMessagesUseCase{repo: repo}
}

func (uc *ListMessagesUseCase) Do(ctx context.Context, userID uuid.UUID) ([]*message.Message, error) {
	return uc.repo.ListMessages(ctx, userID)
}
