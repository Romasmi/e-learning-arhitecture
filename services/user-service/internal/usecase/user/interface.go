package user

import (
	"context"

	"github.com/Romasmi/e-learning-arhitecture/user-service/internal/domain/user"
	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(ctx context.Context, u *user.User) (*user.User, error)
	GetUserById(ctx context.Context, userId uuid.UUID) (*user.User, error)
	UpdateUser(ctx context.Context, u *user.User) (*user.User, error)
	DeleteUser(ctx context.Context, userId uuid.UUID) error
}

type EventProducer interface {
	EmitUserCreated(ctx context.Context, u *user.User) error
}
