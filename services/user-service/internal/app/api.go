package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	api "github.com/Romasmi/e-learning-arhitecture/gen/go/user"
	grpcint "github.com/Romasmi/e-learning-arhitecture/user-service/internal/transport/grpc"
	"github.com/Romasmi/e-learning-arhitecture/user-service/internal/transport/gw"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Api struct {
	App        *App
	grpcServer *grpc.Server
	gwServer   *http.Server
}

func NewApi(app *App) *Api {
	return &Api{App: app}
}

func (a *Api) Run() error {
	grpcAddr := fmt.Sprintf(":%d", a.App.Config.Server.GRPCPort)
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	a.grpcServer = grpc.NewServer()
	userHandler := grpcint.NewUserHandler(a.App)
	api.RegisterUserServiceServer(a.grpcServer, userHandler)
	reflection.Register(a.grpcServer)

	go func() {
		slog.Info("Starting gRPC server", "addr", grpcAddr)
		if err := a.grpcServer.Serve(lis); err != nil {
			slog.Error("gRPC server error", "error", err)
		}
	}()

	var errGw error
	a.gwServer, errGw = gw.NewGatewayServer(a.App.DbConn, grpcAddr, a.App.Config.Server.Port)
	if errGw != nil {
		return fmt.Errorf("failed to create gateway server: %w", errGw)
	}

	go func() {
		slog.Info("Starting HTTP gateway", "addr", a.gwServer.Addr)
		if err := a.gwServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP gateway error", "error", err)
		}
	}()

	return nil
}

func (a *Api) Shutdown(ctx context.Context) error {
	if a.grpcServer != nil {
		slog.Info("Shutting down gRPC server...")
		a.grpcServer.GracefulStop()
	}
	if a.gwServer != nil {
		slog.Info("Shutting down HTTP gateway...")
		if err := a.gwServer.Shutdown(ctx); err != nil {
			slog.Error("HTTP gateway shutdown error", "error", err)
		}
	}
	return a.App.Shutdown(ctx)
}
