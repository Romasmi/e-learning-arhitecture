package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	grpcint "github.com/Romasmi/e-learning-arhitecture/notification-service/internal/interface/grpc"
	httpint "github.com/Romasmi/e-learning-arhitecture/notification-service/internal/interface/http"
)

type Api struct {
	*App
}

func NewApi(app *App) *Api {
	return &Api{App: app}
}

func (a *Api) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// gRPC Server
	grpcServer := grpcint.NewServer(a.App)

	grpcAddr := fmt.Sprintf(":%d", a.Cfg.Server.GRPCPort)
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	go func() {
		slog.Info("Starting gRPC server", "addr", grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("gRPC server error", "error", err)
		}
	}()

	// gRPC Gateway
	gwServer, err := httpint.NewGatewayServer(a.App, grpcAddr, a.Cfg.Server.HTTPPort)
	if err != nil {
		return fmt.Errorf("failed to create gateway server: %w", err)
	}

	go func() {
		slog.Info("Starting HTTP gateway", "addr", gwServer.Addr)
		if err := gwServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP gateway error", "error", err)
		}
	}()

	<-ctx.Done()
	slog.Info("Shutting down API...")
	grpcServer.GracefulStop()
	if err := gwServer.Shutdown(context.Background()); err != nil {
		slog.Error("HTTP gateway shutdown error", "error", err)
	}
	return nil
}
