package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/app"
	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/infra/database"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	basePath := os.Getenv("APP_BASE_PATH")
	if basePath == "" {
		basePath = "."
	}

	appInstance, err := app.NewApp(basePath)
	if err != nil {
		slog.Error("error while app init", "error", err)
		os.Exit(1)
	}
	api := app.NewApi(appInstance)

	migrationsPath := filepath.Join(basePath, "migrations")
	absPath, _ := filepath.Abs(migrationsPath)
	m, err := migrate.New(
		"file://"+absPath,
		database.GetDbUrl(&appInstance.Config.Db),
	)
	if m == nil || err != nil {
		slog.Error("unable to create migrations driver", "error", err)
		os.Exit(1)
	}
	defer func() {
		if sourceErr, dbErr := m.Close(); sourceErr != nil || dbErr != nil {
			slog.Error("Error closing migration driver", "source_err", sourceErr, "db_err", dbErr)
		}
	}()

	err = m.Up()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("error while running up migrations", "error", err)
		os.Exit(1)
	}

	go func() {
		if err := api.Run(); err != nil {
			slog.Error("server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := api.Shutdown(ctx); err != nil {
		slog.Error("Error during shutdown", "error", err)
	}

	slog.Info("Server stopped")
}
