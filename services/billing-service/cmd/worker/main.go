package main

import (
	"log/slog"
	"os"

	"github.com/Romasmi/e-learning-arhitecture/billing-service/internal/app"
	"github.com/Romasmi/e-learning-arhitecture/billing-service/internal/config"
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

	application, err := app.NewApp(cfg)
	if err != nil {
		slog.Error("Failed to initialize app", "error", err)
		os.Exit(1)
	}
	defer application.Close()

	worker := app.NewWorker(application)
	if err := worker.Run(); err != nil {
		slog.Error("Worker error", "error", err)
		os.Exit(1)
	}
}
