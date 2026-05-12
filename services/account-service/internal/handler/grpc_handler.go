package handler

import (
	"context"

	accountpb "github.com/Romasmi/e-learning-arhitecture/gen/go/account"
	"github.com/elearning/account-service/internal/domain"
	"github.com/elearning/account-service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCHandler struct {
	accountpb.UnimplementedAccountServiceServer
	usecase *usecase.AccountUsecase
}

func NewGRPCHandler(u *usecase.AccountUsecase) *GRPCHandler {
	return &GRPCHandler{usecase: u}
}

func (h *GRPCHandler) CreateAccount(ctx context.Context, req *accountpb.CreateAccountRequest) (*accountpb.CreateAccountResponse, error) {
	account, err := h.usecase.CreateAccount(ctx, req.PortalId, req.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create account: %v", err)
	}

	return &accountpb.CreateAccountResponse{
		Account: mapAccountToProto(account),
	}, nil
}

func (h *GRPCHandler) ArchiveAccount(ctx context.Context, req *accountpb.ArchiveAccountRequest) (*accountpb.ArchiveAccountResponse, error) {
	archived, err := h.usecase.ArchiveAccount(ctx, req.AccountId)
	if err != nil {
		if err == domain.ErrAccountNotFound {
			return nil, status.Error(codes.NotFound, "account not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to archive account: %v", err)
	}

	return &accountpb.ArchiveAccountResponse{Archived: archived}, nil
}

func (h *GRPCHandler) CreateAdmin(ctx context.Context, req *accountpb.CreateAdminRequest) (*accountpb.CreateAdminResponse, error) {
	admin, err := h.usecase.CreateAdmin(ctx, req.AccountId, req.Email, req.Password)
	if err != nil {
		if err == domain.ErrAccountNotFound {
			return nil, status.Error(codes.NotFound, "account not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to create admin: %v", err)
	}

	return &accountpb.CreateAdminResponse{
		Admin: mapAdminToProto(admin),
	}, nil
}

func (h *GRPCHandler) GetAccount(ctx context.Context, req *accountpb.GetAccountRequest) (*accountpb.GetAccountResponse, error) {
	account, err := h.usecase.GetAccount(ctx, req.AccountId)
	if err != nil {
		if err == domain.ErrAccountNotFound {
			return nil, status.Error(codes.NotFound, "account not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get account: %v", err)
	}

	return &accountpb.GetAccountResponse{
		Account: mapAccountToProto(account),
	}, nil
}

func (h *GRPCHandler) ListAccounts(ctx context.Context, req *accountpb.ListAccountsRequest) (*accountpb.ListAccountsResponse, error) {
	accounts, err := h.usecase.ListAccounts(ctx, req.PortalId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list accounts: %v", err)
	}

	pbAccounts := make([]*accountpb.Account, len(accounts))
	for i, a := range accounts {
		pbAccounts[i] = mapAccountToProto(a)
	}

	return &accountpb.ListAccountsResponse{Accounts: pbAccounts}, nil
}

func mapAccountToProto(a *domain.Account) *accountpb.Account {
	var status accountpb.AccountStatus
	switch a.Status {
	case domain.AccountStatusActive:
		status = accountpb.AccountStatus_ACTIVE
	case domain.AccountStatusArchived:
		status = accountpb.AccountStatus_ARCHIVED
	default:
		status = accountpb.AccountStatus_ACCOUNT_STATUS_UNSPECIFIED
	}

	return &accountpb.Account{
		Id:        a.ID,
		PortalId:  a.PortalID,
		Name:      a.Name,
		Status:    status,
		CreatedAt: timestamppb.New(a.CreatedAt),
	}
}

func mapAdminToProto(a *domain.Admin) *accountpb.Admin {
	return &accountpb.Admin{
		Id:        a.ID,
		AccountId: a.AccountID,
		Email:     a.Email,
		CreatedAt: timestamppb.New(a.CreatedAt),
	}
}
