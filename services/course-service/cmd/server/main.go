package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	coursepb "github.com/Romasmi/e-learning-arhitecture/gen/go/course"
	"github.com/elearning/course-service/internal/handler"
	"github.com/elearning/course-service/internal/metrics"
	"github.com/elearning/course-service/internal/repository/postgres"
	"github.com/elearning/course-service/internal/usecase"
	"github.com/elearning/course-service/pkg/kafka"
	"github.com/elearning/course-service/pkg/logger"
	pkgpostgres "github.com/elearning/course-service/pkg/postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	kgo "github.com/segmentio/kafka-go"
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

	// Kafka Producer
	producer := kafka.NewProducer(kafkaBrokers, "course.events")
	defer producer.Close()

	// Layers
	repo := postgres.NewCourseRepository(pool)
	uc := usecase.NewCourseUsecase(repo, producer)
	h := handler.NewGRPCHandler(uc)

	// Metrics
	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = "9090"
	}
	metricsServer := &http.Server{
		Addr:    ":" + metricsPort,
		Handler: promhttp.Handler(),
	}

	go func() {
		log.Info("Starting metrics server", zap.String("addr", metricsServer.Addr))
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("metrics server failed", zap.Error(err))
		}
	}()

	// Kafka Consumer for media events
	consumer := kafka.NewConsumer(kafkaBrokers, "media.events", "course-service-group")
	go consumer.Consume(ctx, func(msg kgo.Message) error {
		var event struct {
			EventType string `json:"event_type"`
			AssetID   string `json:"asset_id"`
			Status    string `json:"status"`
		}
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			return err
		}

		if event.EventType == "VideoProcessed" {
			log.Info("Handling VideoProcessed event", zap.String("asset_id", event.AssetID))
			return uc.HandleVideoProcessed(ctx, event.AssetID, "READY")
		} else if event.EventType == "VideoProcessingFailed" {
			log.Info("Handling VideoProcessingFailed event", zap.String("asset_id", event.AssetID))
			return uc.HandleVideoProcessed(ctx, event.AssetID, "FAILED")
		}

		return nil
	})
	defer consumer.Close()

	// gRPC
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(metrics.GrpcUnaryInterceptor),
	)
	coursepb.RegisterCourseServiceServer(grpcServer, h)
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
	err = coursepb.RegisterCourseServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	if err != nil {
		log.Fatal("failed to register gateway", zap.Error(err))
	}

	gwAddr := ":" + gwPort
	gwServer := &http.Server{
		Addr:    gwAddr,
		Handler: metrics.HttpMiddleware(mux),
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
	if err := metricsServer.Shutdown(shutdownCtx); err != nil {
		log.Error("Metrics shutdown failed", zap.Error(err))
	}
}
