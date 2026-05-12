package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/Romasmi/e-learning-arhitecture/notification-service/internal/app"
	"github.com/Romasmi/e-learning-arhitecture/notification-service/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	basePath := os.Getenv("APP_BASE_PATH")
	if basePath == "" {
		basePath = "."
	}

	cfg, err := config.LoadConfig(basePath)
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// Run migrations
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port, cfg.Db.Name)
	migrationsPath := filepath.Join(basePath, "migrations")
	absPath, _ := filepath.Abs(migrationsPath)
	m, err := migrate.New("file://"+absPath, dbUrl)
	if err != nil {
		slog.Error("Failed to create migration driver", "error", err)
		os.Exit(1)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}
	m.Close()

	application, err := app.NewApp(cfg)
	if err != nil {
		slog.Error("Failed to initialize app", "error", err)
		os.Exit(1)
	}
	defer application.Close()

	api := app.NewApi(application)
	if err := api.Run(); err != nil {
		slog.Error("API error", "error", err)
		os.Exit(1)
	}
}
