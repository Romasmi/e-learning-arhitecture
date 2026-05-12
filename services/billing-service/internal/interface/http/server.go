package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Romasmi/e-learning-arhitecture/billing-service/internal/interface/http/middleware"
	api "github.com/Romasmi/e-learning-arhitecture/gen/go/billing"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ReadyChecker interface {
	Ping(ctx context.Context) error
}

func NewGatewayServer(checker ReadyChecker, grpcAddr string, httpPort uint) (*http.Server, error) {
	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := api.RegisterBillingServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	if err != nil {
		return nil, err
	}

	mainMux := http.NewServeMux()
	mainMux.Handle("/", mux)
	mainMux.Handle("/metrics", promhttp.Handler())

	mainMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mainMux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		if err := checker.Ping(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("UNREADY"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("READY"))
	})

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: middleware.MetricsMiddleware(mainMux),
	}, nil
}
