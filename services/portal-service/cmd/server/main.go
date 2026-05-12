package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	portalpb "github.com/Romasmi/e-learning-arhitecture/gen/go/portal"
	"github.com/elearning/portal-service/internal/handler"
	"github.com/elearning/portal-service/internal/repository/postgres"
	"github.com/elearning/portal-service/internal/usecase"
	"github.com/elearning/portal-service/pkg/kafka"
	"github.com/elearning/portal-service/pkg/logger"
	pkgpostgres "github.com/elearning/portal-service/pkg/postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	_ = godotenv.Load()

	grpcPort := os.Getenv("PORT")
	if grpcPort == "" {
		grpcPort = "8080"
	}
	gwPort := os.Getenv("GW_PORT")
	if gwPort == "" {
		gwPort = "8000"
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)
	}
	kafkaBrokersStr := os.Getenv("KAFKA_BROKERS")
	kafkaBrokers := strings.Split(kafkaBrokersStr, ",")

	log, err := logger.New("info")
	if err != nil {
		fmt.Printf("failed to init logger: %v\n", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Migrations
	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		log.Fatal("failed to create migration instance", zap.Error(err))
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("failed to run migrations", zap.Error(err))
	}
	m.Close()

	// DB
	pool, err := pkgpostgres.NewPool(ctx, dbURL)
	if err != nil {
		log.Fatal("failed to connect to db", zap.Error(err))
	}
	defer pool.Close()

	// Kafka
	producer := kafka.NewProducer(kafkaBrokers, "portal.events")
	defer producer.Close()

	// Layers
	repo := postgres.NewPortalRepository(pool)
	uc := usecase.NewPortalUsecase(repo, producer)
	h := handler.NewGRPCHandler(uc)

	// gRPC
	grpcServer := grpc.NewServer()
	portalpb.RegisterPortalServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	grpcAddr := ":" + grpcPort
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal("failed to listen", zap.Error(err))
	}

	go func() {
		log.Info("Starting gRPC server", zap.String("addr", grpcAddr))
		if err := grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			log.Error("gRPC server failed", zap.Error(err))
		}
	}()

	// Gateway
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = portalpb.RegisterPortalServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	if err != nil {
		log.Fatal("failed to register gateway", zap.Error(err))
	}

	gwAddr := ":" + gwPort
	gwServer := &http.Server{
		Addr:    gwAddr,
		Handler: mux,
	}

	go func() {
		log.Info("Starting Gateway server", zap.String("addr", gwAddr))
		if err := gwServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Gateway server failed", zap.Error(err))
		}
	}()

	<-ctx.Done()
	log.Info("Shutting down...")

	grpcServer.GracefulStop()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := gwServer.Shutdown(shutdownCtx); err != nil {
		log.Error("Gateway shutdown failed", zap.Error(err))
	}
}
