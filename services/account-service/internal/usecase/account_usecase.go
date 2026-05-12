package usecase

import (
	"context"
	"encoding/json"
	"time"

	authapi "github.com/Romasmi/e-learning-arhitecture/gen/go/auth"
	"github.com/elearning/account-service/internal/domain"
	"github.com/elearning/account-service/pkg/kafka"
	"github.com/google/uuid"
)

type AccountUsecase struct {
	repo       domain.AccountRepository
	producer   *kafka.Producer
	authClient authapi.AuthServiceClient
}

func NewAccountUsecase(repo domain.AccountRepository, producer *kafka.Producer, authClient authapi.AuthServiceClient) *AccountUsecase {
	return &AccountUsecase{
		repo:       repo,
		producer:   producer,
		authClient: authClient,
	}
}

func (u *AccountUsecase) CreateAccount(ctx context.Context, portalID, name string) (*domain.Account, error) {
	account := &domain.Account{
		ID:        uuid.New().String(),
		PortalID:  portalID,
		Name:      name,
		Status:    domain.AccountStatusActive,
		CreatedAt: time.Now(),
	}

	if err := u.repo.CreateAccount(ctx, account); err != nil {
		return nil, err
	}

	payload, _ := json.Marshal(account)
	u.producer.PublishAsync(kafka.Event{
		EventType:  "AccountCreated",
		AccountID:  account.ID,
		Payload:    payload,
		OccurredAt: time.Now(),
	})

	return account, nil
}

func (u *AccountUsecase) ArchiveAccount(ctx context.Context, accountID string) (bool, error) {
	if err := u.repo.UpdateAccountStatus(ctx, accountID, domain.AccountStatusArchived); err != nil {
		return false, err
	}

	u.producer.PublishAsync(kafka.Event{
		EventType:  "AccountArchived",
		AccountID:  accountID,
		OccurredAt: time.Now(),
	})

	return true, nil
}

func (u *AccountUsecase) CreateAdmin(ctx context.Context, accountID, email, password string) (*domain.Admin, error) {
	account, err := u.repo.GetAccountByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	resp, err := u.authClient.Register(ctx, &authapi.RegisterRequest{
		Email:    email,
		Password: password,
		PortalId: account.PortalID,
		Role:     "ADMIN",
	})
	if err != nil {
		return nil, err
	}

	admin := &domain.Admin{
		ID:        uuid.New().String(),
		AccountID: accountID,
		UserID:    resp.UserId,
		Email:     email,
		CreatedAt: time.Now(),
	}

	if err := u.repo.CreateAdmin(ctx, admin); err != nil {
		return nil, err
	}

	payload, _ := json.Marshal(admin)
	u.producer.PublishAsync(kafka.Event{
		EventType:  "AdminCreated",
		AccountID:  accountID,
		Payload:    payload,
		OccurredAt: time.Now(),
	})

	return admin, nil
}

func (u *AccountUsecase) GetAccount(ctx context.Context, accountID string) (*domain.Account, error) {
	return u.repo.GetAccountByID(ctx, accountID)
}

func (u *AccountUsecase) ListAccounts(ctx context.Context, portalID string) ([]*domain.Account, error) {
	return u.repo.ListAccounts(ctx, portalID)
}
