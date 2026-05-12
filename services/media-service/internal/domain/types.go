package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type AssetType string

const (
	AssetTypeVideo AssetType = "VIDEO"
	AssetTypePDF   AssetType = "PDF"
	AssetTypeImage AssetType = "IMAGE"
)

type AssetStatus string

const (
	AssetStatusPending AssetStatus = "PENDING"
	AssetStatusReady   AssetStatus = "READY"
	AssetStatusFailed  AssetStatus = "FAILED"
)

type Asset struct {
	ID        uuid.UUID
	LessonID  uuid.UUID
	Type      AssetType
	Status    AssetStatus
	RawURL    string
	CDNURLs   []string
	JobID     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrAssetNotFound = errors.New("asset not found")
)
