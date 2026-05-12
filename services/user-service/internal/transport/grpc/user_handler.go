package grpc

import (
	"context"

	api "github.com/Romasmi/e-learning-arhitecture/gen/go/user"
	"github.com/Romasmi/e-learning-arhitecture/user-service/internal/domain/user"
	"github.com/Romasmi/e-learning-arhitecture/user-service/internal/usecase"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserHandler struct {
	api.UnimplementedUserServiceServer
	app interface {
		GetHandler(id usecase.UseCaseID) usecase.Handler
	}
}

func NewUserHandler(app interface {
	GetHandler(id usecase.UseCaseID) usecase.Handler
}) *UserHandler {
	return &UserHandler{
		app: app,
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *api.CreateUserRequest) (*api.User, error) {
	u := &user.User{
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
	}

	handler := h.app.GetHandler(usecase.UseCaseCreateUser)
	resp, err := handler.Do(ctx, u)
	if err != nil {
		return nil, err
	}

	newUser := resp.(*user.User)
	return &api.User{
		Id:        newUser.ID.String(),
		Username:  newUser.Username,
		Email:     newUser.Email,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Phone:     newUser.Phone,
		Password:  newUser.Password,
		CreatedAt: timestamppb.New(newUser.CreatedAt),
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *api.GetUserRequest) (*api.User, error) {
	id, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	handler := h.app.GetHandler(usecase.UseCaseGetUser)
	resp, err := handler.Do(ctx, id)
	if err != nil {
		return nil, err
	}

	u := resp.(*user.User)
	return &api.User{
		Id:        u.ID.String(),
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Phone:     u.Phone,
		Password:  u.Password,
		CreatedAt: timestamppb.New(u.CreatedAt),
	}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *api.UpdateUserRequest) (*api.User, error) {
	id, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	u := &user.User{
		ID: id,
	}
	if req.Email != nil {
		u.Email = *req.Email
	}
	if req.FirstName != nil {
		u.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		u.LastName = *req.LastName
	}
	if req.Phone != nil {
		u.Phone = *req.Phone
	}

	handler := h.app.GetHandler(usecase.UseCaseUpdateUser)
	resp, err := handler.Do(ctx, u)
	if err != nil {
		return nil, err
	}

	updatedUser := resp.(*user.User)
	return &api.User{
		Id:        updatedUser.ID.String(),
		Username:  updatedUser.Username,
		Email:     updatedUser.Email,
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		Phone:     updatedUser.Phone,
		Password:  updatedUser.Password,
		CreatedAt: timestamppb.New(updatedUser.CreatedAt),
	}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *api.DeleteUserRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	handler := h.app.GetHandler(usecase.UseCaseDeleteUser)
	_, err = handler.Do(ctx, id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
