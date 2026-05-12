package billing

import (
	"context"

	"github.com/Romasmi/e-learning-arhitecture/billing-service/internal/domain/account"
	"github.com/google/uuid"
)

type GetAccountUseCase struct {
	repo Repository
}

func NewGetAccountUseCase(repo Repository) *GetAccountUseCase {
	return &GetAccountUseCase{repo: repo}
}

func (uc *GetAccountUseCase) Do(ctx context.Context, userID uuid.UUID) (*account.Account, error) {
	return uc.repo.GetAccount(ctx, userID)
}
