package grpc

import (
	"context"
	"errors"

	"github.com/Romasmi/e-learning-arhitecture/billing-service/internal/domain/account"
	"github.com/Romasmi/e-learning-arhitecture/billing-service/internal/usecase"
	billinguc "github.com/Romasmi/e-learning-arhitecture/billing-service/internal/usecase/billing"
	api "github.com/Romasmi/e-learning-arhitecture/gen/go/billing"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BillingHandler struct {
	api.UnimplementedBillingServiceServer
	app interface {
		GetHandler(id usecase.UseCaseID) usecase.Handler
	}
}

func NewBillingHandler(app interface {
	GetHandler(id usecase.UseCaseID) usecase.Handler
}) *BillingHandler {
	return &BillingHandler{app: app}
}

func (h *BillingHandler) GetAccount(ctx context.Context, req *api.GetAccountRequest) (*api.Account, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	handler := h.app.GetHandler(usecase.UseCaseGetAccount)
	resp, err := handler.Do(ctx, userID)
	if err != nil {
		if errors.Is(err, account.ErrAccountNotFound) {
			return nil, status.Errorf(codes.NotFound, "account not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get account: %v", err)
	}

	acc := resp.(*account.Account)
	return &api.Account{
		UserId:  acc.UserID.String(),
		Balance: acc.Balance,
	}, nil
}

func (h *BillingHandler) TopUp(ctx context.Context, req *api.TopUpRequest) (*api.Account, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	handler := h.app.GetHandler(usecase.UseCaseTopUp)
	resp, err := handler.Do(ctx, billinguc.TopUpInput{
		UserID: userID,
		Amount: req.Amount,
	})
	if err != nil {
		if errors.Is(err, account.ErrAccountNotFound) {
			return nil, status.Errorf(codes.NotFound, "account not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to top up: %v", err)
	}

	acc := resp.(*account.Account)
	return &api.Account{
		UserId:  acc.UserID.String(),
		Balance: acc.Balance,
	}, nil
}

func (h *BillingHandler) Withdraw(ctx context.Context, req *api.WithdrawRequest) (*api.WithdrawResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	handler := h.app.GetHandler(usecase.UseCaseWithdraw)
	resp, err := handler.Do(ctx, billinguc.WithdrawInput{
		UserID:         userID,
		Amount:         req.Amount,
		IdempotencyKey: req.IdempotencyKey,
	})
	if err != nil {
		if errors.Is(err, account.ErrAccountNotFound) {
			return nil, status.Errorf(codes.NotFound, "account not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to withdraw: %v", err)
	}

	output := resp.(*billinguc.WithdrawOutput)
	return &api.WithdrawResponse{
		Success: output.Success,
		Reason:  output.Reason,
		UpdatedAccount: &api.Account{
			UserId:  output.UpdatedAccount.UserID.String(),
			Balance: output.UpdatedAccount.Balance,
		},
	}, nil
}
