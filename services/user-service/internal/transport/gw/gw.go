package gw

import (
	"context"
	"fmt"
	"net/http"

	api "github.com/Romasmi/e-learning-arhitecture/gen/go/user"
	"github.com/Romasmi/e-learning-arhitecture/user-service/internal/middleware"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ReadyChecker interface {
	Ping() error
}

func NewGatewayServer(checker ReadyChecker, grpcAddr string, httpPort uint) (*http.Server, error) {
	ctx := context.Background()
	gwMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := api.RegisterUserServiceHandlerFromEndpoint(ctx, gwMux, grpcAddr, opts)
	if err != nil {
		return nil, err
	}

	r := mux.NewRouter()

	// Apply AuthMiddleware to /user/{userId} routes
	userRouter := r.PathPrefix("/user/{userId}").Subrouter()
	userRouter.Use(middleware.AuthMiddleware)
	userRouter.Path("").Handler(gwMux)
	userRouter.Path("/").Handler(gwMux)

	// Other routes
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
	r.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		if err := checker.Ping(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("UNREADY"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("READY"))
	})

	// Fallback for other /user routes (like POST /user)
	r.PathPrefix("/").Handler(gwMux)

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: middleware.MetricsMiddleware(r),
	}, nil
}
