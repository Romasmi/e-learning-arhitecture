package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/elearning/portal-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type portalRepository struct {
	pool *pgxpool.Pool
}

func NewPortalRepository(pool *pgxpool.Pool) domain.PortalRepository {
	return &portalRepository{pool: pool}
}

func (r *portalRepository) Create(ctx context.Context, p *domain.Portal) error {
	lmsConfigJSON, err := json.Marshal(p.LMSConfig)
	if err != nil {
		return fmt.Errorf("marshal lms config: %w", err)
	}

	query := `
		INSERT INTO portals (id, code, name, status, lms_config, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = r.pool.Exec(ctx, query, p.ID, p.Code, p.Name, p.Status, lmsConfigJSON, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrDuplicateCode
		}
		return fmt.Errorf("exec insert: %w", err)
	}
	return nil
}

func (r *portalRepository) GetByID(ctx context.Context, id string) (*domain.Portal, error) {
	query := `SELECT id, code, name, status, lms_config, created_at, updated_at FROM portals WHERE id = $1`
	var p domain.Portal
	var lmsConfigJSON []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(&p.ID, &p.Code, &p.Name, &p.Status, &lmsConfigJSON, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPortalNotFound
		}
		return nil, fmt.Errorf("query row: %w", err)
	}

	if err := json.Unmarshal(lmsConfigJSON, &p.LMSConfig); err != nil {
		return nil, fmt.Errorf("unmarshal lms config: %w", err)
	}

	return &p, nil
}

func (r *portalRepository) GetByCode(ctx context.Context, code string) (*domain.Portal, error) {
	query := `SELECT id, code, name, status, lms_config, created_at, updated_at FROM portals WHERE code = $1`
	var p domain.Portal
	var lmsConfigJSON []byte
	err := r.pool.QueryRow(ctx, query, code).Scan(&p.ID, &p.Code, &p.Name, &p.Status, &lmsConfigJSON, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPortalNotFound
		}
		return nil, fmt.Errorf("query row: %w", err)
	}

	if err := json.Unmarshal(lmsConfigJSON, &p.LMSConfig); err != nil {
		return nil, fmt.Errorf("unmarshal lms config: %w", err)
	}

	return &p, nil
}

func (r *portalRepository) UpdateConfig(ctx context.Context, id string, config domain.LMSConfig) (*domain.Portal, error) {
	lmsConfigJSON, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("marshal lms config: %w", err)
	}

	query := `
		UPDATE portals 
		SET lms_config = $1, updated_at = NOW() 
		WHERE id = $2 
		RETURNING id, code, name, status, lms_config, created_at, updated_at
	`
	var p domain.Portal
	var resLmsConfigJSON []byte
	err = r.pool.QueryRow(ctx, query, lmsConfigJSON, id).Scan(&p.ID, &p.Code, &p.Name, &p.Status, &resLmsConfigJSON, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPortalNotFound
		}
		return nil, fmt.Errorf("query row update: %w", err)
	}

	if err := json.Unmarshal(resLmsConfigJSON, &p.LMSConfig); err != nil {
		return nil, fmt.Errorf("unmarshal lms config: %w", err)
	}

	return &p, nil
}

func (r *portalRepository) UpdateStatus(ctx context.Context, id string, status domain.PortalStatus) error {
	query := `UPDATE portals SET status = $1, updated_at = NOW() WHERE id = $2`
	result, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("exec update status: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrPortalNotFound
	}
	return nil
}

func (r *portalRepository) List(ctx context.Context) ([]*domain.Portal, error) {
	query := `SELECT id, code, name, status, lms_config, created_at, updated_at FROM portals`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var portals []*domain.Portal
	for rows.Next() {
		var p domain.Portal
		var lmsConfigJSON []byte
		err := rows.Scan(&p.ID, &p.Code, &p.Name, &p.Status, &lmsConfigJSON, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		if err := json.Unmarshal(lmsConfigJSON, &p.LMSConfig); err != nil {
			return nil, fmt.Errorf("unmarshal lms config: %w", err)
		}
		portals = append(portals, &p)
	}

	return portals, nil
}
