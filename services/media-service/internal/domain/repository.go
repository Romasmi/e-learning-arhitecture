package domain

import (
	"context"

	"github.com/google/uuid"
)

type AssetRepository interface {
	CreateAsset(ctx context.Context, asset *Asset) error
	GetAsset(ctx context.Context, id uuid.UUID) (*Asset, error)
	UpdateAssetStatus(ctx context.Context, id uuid.UUID, status AssetStatus, cdnURLs []string) error
	UpdateAssetJob(ctx context.Context, id uuid.UUID, jobID string) error
	GetAssetByJobID(ctx context.Context, jobID string) (*Asset, error)
}
