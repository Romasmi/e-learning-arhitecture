package grpc

import (
	"github.com/Romasmi/e-learning-arhitecture/billing-service/internal/usecase"
	api "github.com/Romasmi/e-learning-arhitecture/gen/go/billing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewServer(app interface {
	GetHandler(id usecase.UseCaseID) usecase.Handler
}) *grpc.Server {
	grpcServer := grpc.NewServer()

	billingHandler := NewBillingHandler(app)
	api.RegisterBillingServiceServer(grpcServer, billingHandler)

	reflection.Register(grpcServer)

	return grpcServer
}
