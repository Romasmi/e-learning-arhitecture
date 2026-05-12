package app

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/routes"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Api struct {
	App    *App
	router *mux.Router
	server *http.Server
}

func NewApi(app *App) *Api {
	api := &Api{App: app}
	api.init()
	return api
}

func (a *Api) init() {
	router := mux.NewRouter()
	routes.RegisterRoutes(router, a.App)
	a.router = router
}

func (a *Api) Run() error {
	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodOptions,
		http.MethodHead,
	})
	headers := handlers.AllowedHeaders([]string{
		"Content-Type",
		"Authorization",
	})
	origins := handlers.AllowedOrigins([]string{"*"})

	a.server = &http.Server{
		Addr:    ":" + strconv.Itoa(int(a.App.Config.Server.Port)),
		Handler: handlers.CORS(credentials, methods, origins, headers)(a.router),
	}

	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (a *Api) Shutdown(ctx context.Context) error {
	return a.App.Shutdown(ctx)
}
