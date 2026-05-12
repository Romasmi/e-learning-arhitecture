package user

import (
	"context"

	"github.com/google/uuid"
)

type DeleteUserUseCase struct {
	repo Repository
}

func NewDeleteUserUseCase(repo Repository) *DeleteUserUseCase {
	return &DeleteUserUseCase{repo: repo}
}

func (uc *DeleteUserUseCase) Do(ctx context.Context, userId uuid.UUID) (any, error) {
	err := uc.repo.DeleteUser(ctx, userId)
	return nil, err
}
