package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/elearning/account-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type accountRepository struct {
	pool *pgxpool.Pool
}

func NewAccountRepository(pool *pgxpool.Pool) domain.AccountRepository {
	return &accountRepository{pool: pool}
}

func (r *accountRepository) CreateAccount(ctx context.Context, a *domain.Account) error {
	query := `
		INSERT INTO accounts (id, portal_id, name, status, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query, a.ID, a.PortalID, a.Name, a.Status, a.CreatedAt)
	if err != nil {
		return fmt.Errorf("exec insert account: %w", err)
	}
	return nil
}

func (r *accountRepository) GetAccountByID(ctx context.Context, id string) (*domain.Account, error) {
	query := `SELECT id, portal_id, name, status, created_at FROM accounts WHERE id = $1`
	var a domain.Account
	err := r.pool.QueryRow(ctx, query, id).Scan(&a.ID, &a.PortalID, &a.Name, &a.Status, &a.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAccountNotFound
		}
		return nil, fmt.Errorf("query row account: %w", err)
	}
	return &a, nil
}

func (r *accountRepository) ListAccounts(ctx context.Context, portalID string) ([]*domain.Account, error) {
	query := `SELECT id, portal_id, name, status, created_at FROM accounts WHERE portal_id = $1`
	rows, err := r.pool.Query(ctx, query, portalID)
	if err != nil {
		return nil, fmt.Errorf("query accounts: %w", err)
	}
	defer rows.Close()

	var accounts []*domain.Account
	for rows.Next() {
		var a domain.Account
		err := rows.Scan(&a.ID, &a.PortalID, &a.Name, &a.Status, &a.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan account: %w", err)
		}
		accounts = append(accounts, &a)
	}
	return accounts, nil
}

func (r *accountRepository) UpdateAccountStatus(ctx context.Context, id string, status domain.AccountStatus) error {
	query := `UPDATE accounts SET status = $1 WHERE id = $2`
	result, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("exec update account status: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrAccountNotFound
	}
	return nil
}

func (r *accountRepository) CreateAdmin(ctx context.Context, a *domain.Admin) error {
	query := `
		INSERT INTO admins (id, account_id, user_id, email, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query, a.ID, a.AccountID, a.UserID, a.Email, a.CreatedAt)
	if err != nil {
		return fmt.Errorf("exec insert admin: %w", err)
	}
	return nil
}
