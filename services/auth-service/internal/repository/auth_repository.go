package repository

import (
	"context"

	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/domain/auth"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	GetByLogin(ctx context.Context, login string) (*auth.Auth, error)
	Create(ctx context.Context, a *auth.Auth) error
	LogAction(ctx context.Context, log *auth.AuthLog) error
}

type authRepository struct {
	db *pgxpool.Pool
}

func CreateAuthRepository(db *pgxpool.Pool) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) GetByLogin(ctx context.Context, login string) (*auth.Auth, error) {
	query := `SELECT id, user_id, login, password_hash, portal_id, role, created_at, updated_at FROM auth WHERE login = $1`
	var a auth.Auth
	err := r.db.QueryRow(ctx, query, login).Scan(
		&a.ID, &a.UserID, &a.Login, &a.PasswordHash, &a.PortalID, &a.Role, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *authRepository) Create(ctx context.Context, a *auth.Auth) error {
	query := `INSERT INTO auth (user_id, login, password_hash, portal_id, role) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, a.UserID, a.Login, a.PasswordHash, a.PortalID, a.Role).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func (r *authRepository) LogAction(ctx context.Context, log *auth.AuthLog) error {
	query := `INSERT INTO auth_logs (user_id, login, action, ip_address) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, log.UserID, log.Login, log.Action, log.IPAddress).Scan(&log.ID, &log.CreatedAt)
}
