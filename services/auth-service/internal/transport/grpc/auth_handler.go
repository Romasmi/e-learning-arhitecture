package grpc

import (
	"context"

	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/services"
	authapi "github.com/Romasmi/e-learning-arhitecture/gen/go/auth"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthHandler struct {
	authapi.UnimplementedAuthServiceServer
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(ctx context.Context, req *authapi.LoginRequest) (*authapi.LoginResponse, error) {
	result, err := h.authService.Login(ctx, req.Email, req.Password, req.PortalId, "")
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &authapi.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresAt:    timestamppb.New(result.ExpiresAt),
		User: &authapi.UserInfo{
			UserId:   result.User.UserID,
			Email:    result.User.Email,
			PortalId: result.User.PortalID,
			Role:     result.User.Role,
		},
	}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *authapi.RefreshTokenRequest) (*authapi.RefreshTokenResponse, error) {
	result, err := h.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &authapi.RefreshTokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresAt:    timestamppb.New(result.ExpiresAt),
	}, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *authapi.ValidateTokenRequest) (*authapi.ValidateTokenResponse, error) {
	claims, err := h.authService.Validate(ctx, req.AccessToken)
	if err != nil {
		return &authapi.ValidateTokenResponse{Valid: false}, nil
	}

	return &authapi.ValidateTokenResponse{
		Valid:     true,
		UserId:    claims.UserID,
		PortalId:  claims.PortalID,
		Role:      claims.Role,
		ExpiresAt: timestamppb.New(claims.ExpiresAt.Time),
	}, nil
}

func (h *AuthHandler) ResetPassword(ctx context.Context, req *authapi.ResetPasswordRequest) (*authapi.ResetPasswordResponse, error) {
	sent, err := h.authService.ResetPassword(ctx, req.Email, req.PortalId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authapi.ResetPasswordResponse{Sent: sent}, nil
}

func (h *AuthHandler) Register(ctx context.Context, req *authapi.RegisterRequest) (*authapi.RegisterResponse, error) {
	var userID uuid.UUID
	var err error
	if req.UserId != "" {
		userID, err = uuid.Parse(req.UserId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid user_id")
		}
	}

	resID, err := h.authService.Register(ctx, userID, req.Email, req.Password, req.PortalId, req.Role)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authapi.RegisterResponse{
		Success: true,
		UserId:  resID.String(),
	}, nil
}
