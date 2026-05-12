package gw

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/handlers/http/auth_handler"
	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/metrics"
	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/services"
	authapi "github.com/Romasmi/e-learning-arhitecture/gen/go/auth"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

func NewGatewayServer(authService services.AuthService, grpcAddr string, httpPort uint) (*http.Server, error) {
	ctx := context.Background()
	gwMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := authapi.RegisterAuthServiceHandlerFromEndpoint(ctx, gwMux, grpcAddr, opts)
	if err != nil {
		return nil, err
	}

	r := mux.NewRouter()

	// Traefik Forward Auth support
	h := auth_handler.NewAuthHandler(authService)
	r.HandleFunc("/auth", h.ValidateHandler).Methods(http.MethodGet, http.MethodHead)

	// Metrics and Health
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Fallback to Gateway Mux
	r.PathPrefix("/").Handler(gwMux)

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: metrics.HttpMiddleware(r),
	}, nil
}
