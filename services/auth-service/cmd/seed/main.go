package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/config"
	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/infra/database"
	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/repository"
	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/services"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/subosito/gotenv"
)

func main() {
	email := flag.String("email", "", "Email for the supervisor")
	password := flag.String("password", "", "Password for the supervisor")
	portalID := flag.String("portal_id", "default", "Portal ID for the supervisor")
	flag.Parse()

	if *email == "" || *password == "" {
		fmt.Println("Usage: go run ./cmd/seed/main.go --email <email> --password <password> [--portal_id <portal_id>]")
		os.Exit(1)
	}

	// Load .env file if it exists
	_ = gotenv.Load()

	basePath := os.Getenv("APP_BASE_PATH")
	if basePath == "" {
		basePath = "."
	}

	cfg, err := config.LoadConfig(basePath)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	dbUrl := database.GetDbUrl(&cfg.Db)
	db, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	authRepo := repository.CreateAuthRepository(db)
	authService := services.CreateAuthService(authRepo, cfg)

	userID := uuid.New()
	role := "supervisor"

	_, err = authService.Upsert(context.Background(), userID, *email, *password, *portalID, role)
	if err != nil {
		slog.Error("failed to upsert supervisor", "error", err)
		os.Exit(1)
	}

	fmt.Printf("Supervisor created: %s / portal: %s\n", *email, *portalID)
}
