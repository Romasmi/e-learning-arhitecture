package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Romasmi/e-learning-arhitecture/notification-service/internal/config"
)

type Worker struct {
	App *App
	Cfg *config.Config
}

func NewWorker(app *App) *Worker {
	return &Worker{
		App: app,
		Cfg: app.Cfg,
	}
}

func (w *Worker) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	for _, consumer := range w.App.Consumers {
		go consumer.Start(ctx)
	}

	<-ctx.Done()
	slog.Info("Shutting down worker...")
	return nil
}
