package grpc

import (
	api "github.com/Romasmi/e-learning-arhitecture/gen/go/notification"
	"github.com/Romasmi/e-learning-arhitecture/notification-service/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewServer(app interface {
	GetHandler(id usecase.UseCaseID) usecase.Handler
}) *grpc.Server {
	grpcServer := grpc.NewServer()

	notificationHandler := NewNotificationHandler(app)
	api.RegisterNotificationServiceServer(grpcServer, notificationHandler)

	reflection.Register(grpcServer)

	return grpcServer
}
