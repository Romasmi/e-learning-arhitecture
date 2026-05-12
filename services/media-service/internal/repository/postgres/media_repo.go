package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/elearning/media-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type assetRepository struct {
	pool *pgxpool.Pool
}

func NewAssetRepository(pool *pgxpool.Pool) domain.AssetRepository {
	return &assetRepository{pool: pool}
}

func (r *assetRepository) CreateAsset(ctx context.Context, a *domain.Asset) error {
	query := `
		INSERT INTO media_assets (id, lesson_id, type, status, raw_url, cdn_urls, job_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.pool.Exec(ctx, query, a.ID, a.LessonID, a.Type, a.Status, a.RawURL, a.CDNURLs, a.JobID, a.CreatedAt, a.UpdatedAt)
	if err != nil {
		return fmt.Errorf("exec insert: %w", err)
	}
	return nil
}

func (r *assetRepository) GetAsset(ctx context.Context, id uuid.UUID) (*domain.Asset, error) {
	query := `SELECT id, lesson_id, type, status, raw_url, cdn_urls, job_id, created_at, updated_at FROM media_assets WHERE id = $1`
	var a domain.Asset
	err := r.pool.QueryRow(ctx, query, id).Scan(&a.ID, &a.LessonID, &a.Type, &a.Status, &a.RawURL, &a.CDNURLs, &a.JobID, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAssetNotFound
		}
		return nil, fmt.Errorf("query row: %w", err)
	}
	return &a, nil
}

func (r *assetRepository) UpdateAssetStatus(ctx context.Context, id uuid.UUID, status domain.AssetStatus, cdnURLs []string) error {
	query := `UPDATE media_assets SET status = $1, cdn_urls = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.pool.Exec(ctx, query, status, cdnURLs, id)
	if err != nil {
		return fmt.Errorf("exec update status: %w", err)
	}
	return nil
}

func (r *assetRepository) UpdateAssetJob(ctx context.Context, id uuid.UUID, jobID string) error {
	query := `UPDATE media_assets SET job_id = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, jobID, id)
	if err != nil {
		return fmt.Errorf("exec update job: %w", err)
	}
	return nil
}

func (r *assetRepository) GetAssetByJobID(ctx context.Context, jobID string) (*domain.Asset, error) {
	query := `SELECT id, lesson_id, type, status, raw_url, cdn_urls, job_id, created_at, updated_at FROM media_assets WHERE job_id = $1`
	var a domain.Asset
	err := r.pool.QueryRow(ctx, query, jobID).Scan(&a.ID, &a.LessonID, &a.Type, &a.Status, &a.RawURL, &a.CDNURLs, &a.JobID, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAssetNotFound
		}
		return nil, fmt.Errorf("query row by job: %w", err)
	}
	return &a, nil
}
