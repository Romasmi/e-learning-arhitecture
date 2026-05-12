package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/elearning/portal-service/internal/domain"
	"github.com/elearning/portal-service/pkg/kafka"
	"github.com/google/uuid"
)

type PortalUsecase struct {
	repo     domain.PortalRepository
	producer *kafka.Producer
}

func NewPortalUsecase(repo domain.PortalRepository, producer *kafka.Producer) *PortalUsecase {
	return &PortalUsecase{
		repo:     repo,
		producer: producer,
	}
}

func (u *PortalUsecase) CreatePortal(ctx context.Context, code, name string, config domain.LMSConfig) (*domain.Portal, error) {
	portal := &domain.Portal{
		ID:        uuid.New().String(),
		Code:      code,
		Name:      name,
		Status:    domain.PortalStatusActive,
		LMSConfig: config,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.repo.Create(ctx, portal); err != nil {
		return nil, err
	}

	payload, _ := json.Marshal(portal)
	u.producer.PublishAsync(kafka.Event{
		EventType:  "PortalCreated",
		PortalID:   portal.ID,
		Payload:    payload,
		OccurredAt: time.Now(),
	})

	return portal, nil
}

func (u *PortalUsecase) GetPortal(ctx context.Context, id string) (*domain.Portal, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *PortalUsecase) UpdatePortalConfig(ctx context.Context, id string, name string, config domain.LMSConfig) (*domain.Portal, error) {
	portal, err := u.repo.Update(ctx, id, name, config)
	if err != nil {
		return nil, err
	}

	payload, _ := json.Marshal(portal.LMSConfig)
	u.producer.PublishAsync(kafka.Event{
		EventType:  "PortalConfigUpdated",
		PortalID:   portal.ID,
		Payload:    payload,
		OccurredAt: time.Now(),
	})

	return portal, nil
}

func (u *PortalUsecase) ArchivePortal(ctx context.Context, id string) error {
	if err := u.repo.UpdateStatus(ctx, id, domain.PortalStatusArchived); err != nil {
		return err
	}

	u.producer.PublishAsync(kafka.Event{
		EventType:  "PortalArchived",
		PortalID:   id,
		OccurredAt: time.Now(),
	})

	return nil
}

func (u *PortalUsecase) ListPortals(ctx context.Context) ([]*domain.Portal, error) {
	return u.repo.List(ctx)
}
