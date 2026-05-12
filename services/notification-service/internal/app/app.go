package app

import (
	"context"
	"fmt"

	"github.com/Romasmi/e-learning-arhitecture/notification-service/internal/config"
	"github.com/Romasmi/e-learning-arhitecture/notification-service/internal/infrastructure/db/postgres"
	"github.com/Romasmi/e-learning-arhitecture/notification-service/internal/interface/kafka"
	"github.com/Romasmi/e-learning-arhitecture/notification-service/internal/usecase"
	notificationuc "github.com/Romasmi/e-learning-arhitecture/notification-service/internal/usecase/notification"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Cfg         *config.Config
	Pool        *pgxpool.Pool
	MessageRepo *postgres.MessageRepository
	Handlers    map[usecase.UseCaseID]usecase.Handler
	Consumers   []kafka.Consumer
}

func NewApp(cfg *config.Config) (*App, error) {
	ctx := context.Background()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port, cfg.Db.Name)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	messageRepo := postgres.NewMessageRepository(pool)

	app := &App{
		Cfg:         cfg,
		Pool:        pool,
		MessageRepo: messageRepo,
		Handlers:    make(map[usecase.UseCaseID]usecase.Handler),
		Consumers:   make([]kafka.Consumer, 0),
	}

	app.registerHandlers()
	app.registerConsumers()

	return app, nil
}

func (a *App) registerHandlers() {
	a.Handlers[usecase.UseCaseListMessages] = usecase.NewHandler(notificationuc.NewListMessagesUseCase(a.MessageRepo))
}

func (a *App) GetHandler(id usecase.UseCaseID) usecase.Handler {
	return a.Handlers[id]
}

func (a *App) registerConsumers() {
	orderConsumer := kafka.NewOrderConsumer(a.Cfg.Kafka.Brokers, a.Cfg.Kafka.Topic, a.MessageRepo)
	a.Consumers = append(a.Consumers, orderConsumer)
}

func (a *App) Close() {
	if a.Pool != nil {
		a.Pool.Close()
	}
	for _, c := range a.Consumers {
		c.Close()
	}
}

func (a *App) GetConfig() *config.Config {
	return a.Cfg
}

func (a *App) Ping(ctx context.Context) error {
	if a.Pool == nil {
		return fmt.Errorf("database pool is not initialized")
	}
	return a.Pool.Ping(ctx)
}
