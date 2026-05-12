package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/config"
	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/infra/database"
)

type App struct {
	DbConn *database.Connection
	Config *config.Config
	server *http.Server
}

func (a *App) GetDB() *database.Connection {
	return a.DbConn
}

func (a *App) GetConfig() *config.Config {
	return a.Config
}

func NewApp(configPath string) (*App, error) {
	app := &App{}
	return app, app.init(configPath)
}

func (a *App) init(configPath string) error {
	envConfig, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("error loading Config: %v\n", err)
	}
	a.Config = envConfig

	dbConn := &database.Connection{Config: &envConfig.Db}
	if err = dbConn.Connect(); err != nil {
		return fmt.Errorf("error connecting to DB: %v\n", err)
	}
	a.DbConn = dbConn

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	var shutdownErr error

	if a.server != nil {
		slog.Info("Shutting down HTTP server...")
		if err := a.server.Shutdown(ctx); err != nil {
			shutdownErr = fmt.Errorf("server shutdown error: %w", err)
			slog.Error("HTTP server shutdown error", "error", err)
		}
	}

	if a.DbConn != nil {
		slog.Info("Closing database connections...")
		a.DbConn.Close()
	}

	slog.Info("Cleanup completed")
	return shutdownErr
}
