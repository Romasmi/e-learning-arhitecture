package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Romasmi/e-learning-arhitecture/user-service/internal/domain/user"
	"github.com/brianvoe/gofakeit/v7"
)

type CreateUserUseCase struct {
	repo           Repository
	producer       EventProducer
	authServiceURL string
}

func NewCreateUserUseCase(repo Repository, producer EventProducer, authServiceURL string) *CreateUserUseCase {
	return &CreateUserUseCase{
		repo:           repo,
		producer:       producer,
		authServiceURL: authServiceURL,
	}
}

func (uc *CreateUserUseCase) Do(ctx context.Context, u *user.User) (*user.User, error) {
	// Generate random password because user doesn't provide it yet
	password := gofakeit.Password(true, true, true, true, false, 12)
	u.Password = password

	newUser, err := uc.repo.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	if err := uc.registerInAuth(ctx, newUser.ID.String(), newUser.Username, password); err != nil {
		// TODO add saga later
		return nil, fmt.Errorf("error registering user in auth service: %w", err)
	}

	newUser.Password = password

	if uc.producer != nil {
		if err := uc.producer.EmitUserCreated(ctx, newUser); err != nil {
			return nil, fmt.Errorf("error emitting user.created event: %w", err)
		}
	}

	return newUser, nil
}

func (uc *CreateUserUseCase) registerInAuth(ctx context.Context, userID, login, password string) error {
	reqBody, _ := json.Marshal(map[string]string{
		"user_id":  userID,
		"login":    login,
		"password": password,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uc.authServiceURL+"/auth/register", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth service returned status %d", resp.StatusCode)
	}

	return nil
}
