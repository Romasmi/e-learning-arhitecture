package user

import (
	"context"

	"github.com/Romasmi/e-learning-arhitecture/user-service/internal/domain/user"
	"github.com/google/uuid"
)

type GetUserUseCase struct {
	repo Repository
}

func NewGetUserUseCase(repo Repository) *GetUserUseCase {
	return &GetUserUseCase{repo: repo}
}

func (uc *GetUserUseCase) Do(ctx context.Context, userId uuid.UUID) (*user.User, error) {
	return uc.repo.GetUserById(ctx, userId)
}
