package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Romasmi/e-learning-arhitecture/user-service/internal/config"
	"github.com/Romasmi/e-learning-arhitecture/user-service/internal/infra/database"
	"github.com/Romasmi/e-learning-arhitecture/user-service/internal/infra/kafka"
	"github.com/Romasmi/e-learning-arhitecture/user-service/internal/repository"
	"github.com/Romasmi/e-learning-arhitecture/user-service/internal/usecase"
	useruc "github.com/Romasmi/e-learning-arhitecture/user-service/internal/usecase/user"
)

type App struct {
	DbConn   *database.Connection
	Config   *config.Config
	Handlers map[usecase.UseCaseID]usecase.Handler
	Producer *kafka.UserProducer
	server   *http.Server
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

	a.Producer = kafka.NewUserProducer(a.Config.Kafka.Brokers, a.Config.Kafka.Topic)

	userRepo := repository.CreateUserRepository(dbConn.DB)

	a.Handlers = make(map[usecase.UseCaseID]usecase.Handler)
	a.registerHandlers(userRepo)

	return nil
}

func (a *App) registerHandlers(userRepo *repository.UserRepository) {
	a.Handlers[usecase.UseCaseCreateUser] = usecase.NewHandler(useruc.NewCreateUserUseCase(userRepo, a.Producer, a.Config.AuthServiceURL))
	a.Handlers[usecase.UseCaseGetUser] = usecase.NewHandler(useruc.NewGetUserUseCase(userRepo))
	a.Handlers[usecase.UseCaseUpdateUser] = usecase.NewHandler(useruc.NewUpdateUserUseCase(userRepo))
	a.Handlers[usecase.UseCaseDeleteUser] = usecase.NewHandler(useruc.NewDeleteUserUseCase(userRepo))
}

func (a *App) GetHandler(id usecase.UseCaseID) usecase.Handler {
	return a.Handlers[id]
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

	if a.Producer != nil {
		slog.Info("Closing Kafka producer...")
		if err := a.Producer.Close(); err != nil {
			slog.Error("Kafka producer close error", "error", err)
		}
	}

	if a.DbConn != nil && a.DbConn.DB != nil {
		slog.Info("Closing database connections...")
		select {
		case <-ctx.Done():
			slog.Warn("Shutdown timeout reached, forcing database close")
		default:
			a.DbConn.DB.Close()
		}
	}

	slog.Info("Cleanup completed")
	return shutdownErr
}
