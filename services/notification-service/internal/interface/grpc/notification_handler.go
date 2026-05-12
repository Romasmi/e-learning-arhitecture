package grpc

import (
	"context"

	api "github.com/Romasmi/e-learning-arhitecture/gen/go/notification"
	"github.com/Romasmi/e-learning-arhitecture/notification-service/internal/domain/message"
	"github.com/Romasmi/e-learning-arhitecture/notification-service/internal/usecase"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotificationHandler struct {
	api.UnimplementedNotificationServiceServer
	app interface {
		GetHandler(id usecase.UseCaseID) usecase.Handler
	}
}

func NewNotificationHandler(app interface {
	GetHandler(id usecase.UseCaseID) usecase.Handler
}) *NotificationHandler {
	return &NotificationHandler{app: app}
}

func (h *NotificationHandler) ListMessages(ctx context.Context, req *api.ListMessagesRequest) (*api.ListMessagesResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	handler := h.app.GetHandler(usecase.UseCaseListMessages)
	resp, err := handler.Do(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list messages: %v", err)
	}

	msgs := resp.([]*message.Message)
	var apiMsgs []*api.Message
	for _, m := range msgs {
		apiMsgs = append(apiMsgs, &api.Message{
			Id:        m.ID.String(),
			UserId:    m.UserID.String(),
			OrderId:   m.OrderID.String(),
			Type:      m.Type,
			Timestamp: m.CreatedAt.String(),
		})
	}

	return &api.ListMessagesResponse{
		Messages: apiMsgs,
	}, nil
}
