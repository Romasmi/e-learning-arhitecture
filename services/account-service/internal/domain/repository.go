package domain

import (
	"context"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, account *Account) error
	GetAccountByID(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context, portalID string) ([]*Account, error)
	UpdateAccountStatus(ctx context.Context, id string, status AccountStatus) error

	CreateAdmin(ctx context.Context, admin *Admin) error
}
