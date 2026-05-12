package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Romasmi/e-learning-arhitecture/billing-service/internal/domain/account"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository struct {
	pool *pgxpool.Pool
}

func NewAccountRepository(pool *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{pool: pool}
}

func (r *AccountRepository) CreateAccount(ctx context.Context, acc *account.Account) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO accounts (user_id, balance, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) ON CONFLICT (user_id) DO NOTHING", acc.UserID, acc.Balance)
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}
	return nil
}

func (r *AccountRepository) GetAccount(ctx context.Context, userID uuid.UUID) (*account.Account, error) {
	acc := &account.Account{}
	err := r.pool.QueryRow(ctx, "SELECT user_id, balance, created_at, updated_at FROM accounts WHERE user_id = $1", userID).Scan(&acc.UserID, &acc.Balance, &acc.CreatedAt, &acc.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, account.ErrAccountNotFound
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return acc, nil
}

func (r *AccountRepository) UpdateBalance(ctx context.Context, userID uuid.UUID, amount int64) (*account.Account, error) {
	acc := &account.Account{}
	err := r.pool.QueryRow(ctx, "UPDATE accounts SET balance = balance + $1, updated_at = NOW() WHERE user_id = $2 RETURNING user_id, balance, created_at, updated_at", amount, userID).Scan(&acc.UserID, &acc.Balance, &acc.CreatedAt, &acc.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, account.ErrAccountNotFound
		}
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}
	return acc, nil
}

func (r *AccountRepository) Withdraw(ctx context.Context, userID uuid.UUID, amount int64) (*account.Account, error) {
	acc := &account.Account{}
	err := r.pool.QueryRow(ctx, "UPDATE accounts SET balance = balance - $1, updated_at = NOW() WHERE user_id = $2 AND balance >= $1 RETURNING user_id, balance, created_at, updated_at", amount, userID).Scan(&acc.UserID, &acc.Balance, &acc.CreatedAt, &acc.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, account.ErrAccountNotFound
		}
		return nil, err
	}
	return acc, nil
}
