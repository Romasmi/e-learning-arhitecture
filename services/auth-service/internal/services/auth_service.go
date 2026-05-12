package services

import (
	"context"
	"fmt"

	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/config"
	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/domain/auth"
	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/repository"
	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/utils/time_utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, login, password, ip string) (string, error)
	Validate(ctx context.Context, token string) (*JWTClaims, error)
	Register(ctx context.Context, userID uuid.UUID, login, password string) error
}

type authService struct {
	repo repository.AuthRepository
	cfg  *config.Config
}

func CreateAuthService(repo repository.AuthRepository, cfg *config.Config) AuthService {
	return &authService{repo: repo, cfg: cfg}
}

func (s *authService) Login(ctx context.Context, login, password, ip string) (string, error) {
	a, err := s.repo.GetByLogin(ctx, login)
	if err != nil {
		s.repo.LogAction(ctx, &auth.AuthLog{
			Login:     login,
			Action:    "LOGIN_FAILED_USER_NOT_FOUND",
			IPAddress: ip,
		})
		return "", fmt.Errorf("invalid login or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(password))
	if err != nil {
		s.repo.LogAction(ctx, &auth.AuthLog{
			UserID:    a.UserID,
			Login:     login,
			Action:    "LOGIN_FAILED_INVALID_PASSWORD",
			IPAddress: ip,
		})
		return "", fmt.Errorf("invalid login or password")
	}

	token, err := GenerateToken(
		a.UserID.String(),
		a.Login,
		s.cfg.Jwt.Secret,
		time_utils.MinutesToDuration(s.cfg.Jwt.Expiration),
	)
	if err != nil {
		return "", err
	}

	s.repo.LogAction(ctx, &auth.AuthLog{
		UserID:    a.UserID,
		Login:     login,
		Action:    "LOGIN_SUCCESS",
		IPAddress: ip,
	})

	return token, nil
}

func (s *authService) Validate(ctx context.Context, token string) (*JWTClaims, error) {
	return ValidateToken(token, s.cfg.Jwt.Secret)
}

func (s *authService) Register(ctx context.Context, userID uuid.UUID, login, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	a := &auth.Auth{
		UserID:       userID,
		Login:        login,
		PasswordHash: string(hashedPassword),
	}

	return s.repo.Create(ctx, a)
}
