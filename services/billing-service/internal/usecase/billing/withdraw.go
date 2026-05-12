package billing

import (
	"context"

	"github.com/Romasmi/e-learning-arhitecture/billing-service/internal/domain/account"
	"github.com/google/uuid"
)

type WithdrawInput struct {
	UserID         uuid.UUID
	Amount         int64
	IdempotencyKey string
}

type WithdrawOutput struct {
	Success        bool
	Reason         string
	UpdatedAccount *account.Account
}

type WithdrawUseCase struct {
	repo Repository
}

func NewWithdrawUseCase(repo Repository) *WithdrawUseCase {
	return &WithdrawUseCase{repo: repo}
}

func (uc *WithdrawUseCase) Do(ctx context.Context, input WithdrawInput) (*WithdrawOutput, error) {
	acc, err := uc.repo.Withdraw(ctx, input.UserID, input.Amount)
	if err != nil {
		currentAcc, fetchErr := uc.repo.GetAccount(ctx, input.UserID)
		if fetchErr == nil {
			return &WithdrawOutput{
				Success:        false,
				Reason:         "insufficient_funds",
				UpdatedAccount: currentAcc,
			}, nil
		}
		return nil, err
	}

	return &WithdrawOutput{
		Success:        true,
		UpdatedAccount: acc,
	}, nil
}
