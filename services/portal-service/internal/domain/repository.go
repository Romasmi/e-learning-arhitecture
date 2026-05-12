package domain

import (
	"context"
)

type PortalRepository interface {
	Create(ctx context.Context, portal *Portal) error
	GetByID(ctx context.Context, id string) (*Portal, error)
	GetByCode(ctx context.Context, code string) (*Portal, error)
	UpdateConfig(ctx context.Context, id string, config LMSConfig) (*Portal, error)
	UpdateStatus(ctx context.Context, id string, status PortalStatus) error
	List(ctx context.Context) ([]*Portal, error)
}
