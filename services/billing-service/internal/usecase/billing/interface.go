package billing

import (
	"context"

	"github.com/Romasmi/e-learning-arhitecture/billing-service/internal/domain/account"
	"github.com/google/uuid"
)

type Repository interface {
	GetAccount(ctx context.Context, userID uuid.UUID) (*account.Account, error)
	UpdateBalance(ctx context.Context, userID uuid.UUID, amount int64) (*account.Account, error)
	Withdraw(ctx context.Context, userID uuid.UUID, amount int64) (*account.Account, error)
	CreateAccount(ctx context.Context, acc *account.Account) error
}
