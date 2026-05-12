package user

import (
	"context"

	"github.com/Romasmi/e-learning-arhitecture/user-service/internal/domain/user"
)

type UpdateUserUseCase struct {
	repo Repository
}

func NewUpdateUserUseCase(repo Repository) *UpdateUserUseCase {
	return &UpdateUserUseCase{repo: repo}
}

func (uc *UpdateUserUseCase) Do(ctx context.Context, u *user.User) (*user.User, error) {
	return uc.repo.UpdateUser(ctx, u)
}
