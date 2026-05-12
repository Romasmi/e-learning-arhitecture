package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/config"
	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/domain/auth"
	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/repository"
	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/utils/time_utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	User         UserInfo
}

type UserInfo struct {
	UserID   string
	Email    string
	PortalID string
	Role     string
}

type AuthService interface {
	Login(ctx context.Context, email, password, portalID, ip string) (*LoginResult, error)
	Validate(ctx context.Context, token string) (*JWTClaims, error)
	Register(ctx context.Context, userID uuid.UUID, login, password, portalID, role string) (uuid.UUID, error)
	Upsert(ctx context.Context, userID uuid.UUID, login, password, portalID, role string) (uuid.UUID, error)
	RefreshToken(ctx context.Context, refreshToken string) (*LoginResult, error)
	ResetPassword(ctx context.Context, email, portalID string) (bool, error)
}

type authService struct {
	repo repository.AuthRepository
	cfg  *config.Config
}

func CreateAuthService(repo repository.AuthRepository, cfg *config.Config) AuthService {
	return &authService{repo: repo, cfg: cfg}
}

func (s *authService) Login(ctx context.Context, email, password, portalID, ip string) (*LoginResult, error) {
	a, err := s.repo.GetByLogin(ctx, email)
	if err != nil {
		s.repo.LogAction(ctx, &auth.AuthLog{
			Login:     email,
			Action:    "LOGIN_FAILED_USER_NOT_FOUND",
			IPAddress: ip,
		})
		return nil, fmt.Errorf("invalid login or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(password))
	if err != nil {
		s.repo.LogAction(ctx, &auth.AuthLog{
			UserID:    a.UserID,
			Login:     email,
			Action:    "LOGIN_FAILED_INVALID_PASSWORD",
			IPAddress: ip,
		})
		return nil, fmt.Errorf("invalid login or password")
	}

	expiration := time_utils.MinutesToDuration(s.cfg.Jwt.Expiration)
	token, err := GenerateToken(
		a.UserID.String(),
		a.Login,
		a.PortalID,
		a.Role,
		s.cfg.Jwt.Secret,
		expiration,
	)
	if err != nil {
		return nil, err
	}

	refreshExpiration := 24 * time.Hour
	refreshToken, err := GenerateToken(
		a.UserID.String(),
		a.Login,
		a.PortalID,
		a.Role,
		s.cfg.Jwt.Secret,
		refreshExpiration,
	)
	if err != nil {
		return nil, err
	}

	s.repo.LogAction(ctx, &auth.AuthLog{
		UserID:    a.UserID,
		Login:     email,
		Action:    "LOGIN_SUCCESS",
		IPAddress: ip,
	})

	return &LoginResult{
		AccessToken:  token,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(expiration),
		User: UserInfo{
			UserID:   a.UserID.String(),
			Email:    a.Login,
			PortalID: a.PortalID,
			Role:     a.Role,
		},
	}, nil
}

func (s *authService) Validate(ctx context.Context, token string) (*JWTClaims, error) {
	return ValidateToken(token, s.cfg.Jwt.Secret)
}

func (s *authService) Register(ctx context.Context, userID uuid.UUID, login, password, portalID, role string) (uuid.UUID, error) {
	if userID == uuid.Nil {
		userID = uuid.New()
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, err
	}

	a := &auth.Auth{
		UserID:       userID,
		Login:        login,
		PasswordHash: string(hashedPassword),
		PortalID:     portalID,
		Role:         role,
	}

	return userID, s.repo.Create(ctx, a)
}

func (s *authService) Upsert(ctx context.Context, userID uuid.UUID, login, password, portalID, role string) (uuid.UUID, error) {
	if userID == uuid.Nil {
		userID = uuid.New()
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, err
	}

	a := &auth.Auth{
		UserID:       userID,
		Login:        login,
		PasswordHash: string(hashedPassword),
		PortalID:     portalID,
		Role:         role,
	}

	if err := s.repo.Upsert(ctx, a); err != nil {
		return uuid.Nil, err
	}
	return a.UserID, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*LoginResult, error) {
	claims, err := ValidateToken(refreshToken, s.cfg.Jwt.Secret)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	expiration := time_utils.MinutesToDuration(s.cfg.Jwt.Expiration)
	accessToken, err := GenerateToken(
		claims.UserID,
		claims.Login,
		claims.PortalID,
		claims.Role,
		s.cfg.Jwt.Secret,
		expiration,
	)
	if err != nil {
		return nil, err
	}

	// For simplicity, reuse the same refresh token or issue a new one
	// Here we issue a new one
	refreshExpiration := 24 * time.Hour
	newRefreshToken, err := GenerateToken(
		claims.UserID,
		claims.Login,
		claims.PortalID,
		claims.Role,
		s.cfg.Jwt.Secret,
		refreshExpiration,
	)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(expiration),
		User: UserInfo{
			UserID:   claims.UserID,
			Email:    claims.Login,
			PortalID: claims.PortalID,
			Role:     claims.Role,
		},
	}, nil
}

func (s *authService) ResetPassword(ctx context.Context, email, portalID string) (bool, error) {
	// Always return true as per contract to avoid email enumeration
	return true, nil
}
