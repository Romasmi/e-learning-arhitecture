package routes

import (
	"net/http"

	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/config"
	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/handlers/http/auth_handler"
	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/repository"
	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/services"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterAuthRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	authRepo := repository.CreateAuthRepository(db)
	authService := services.CreateAuthService(authRepo, cfg)
	authHandler := auth_handler.NewAuthHandler(authService)

	r.HandleFunc("/auth", authHandler.ValidateHandler).Methods(http.MethodGet, http.MethodHead)
	r.HandleFunc("/auth/register", authHandler.RegisterHandler).Methods(http.MethodPost)
	r.HandleFunc("/auth/login", authHandler.LoginHandler).Methods(http.MethodPost)
}
