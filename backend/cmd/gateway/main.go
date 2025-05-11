package main

import (
	pb "University-Selection-Service/pkg/api"
	"University-Selection-Service/pkg/logger"
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			allowedOrigins := map[string]bool{
				"http://localhost:5173": true,
			}
			if allowedOrigins[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
				w.Header().Set("Access-Control-Max-Age", "86400")
				w.Header().Set("Vary", "Origin")
			}
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	ctx := context.Background()
	ctx, _ = logger.NewLogger(ctx)
	log := logger.GetLoggerFromCtx(ctx)

	mux := runtime.NewServeMux()
	muxWithCORS := corsMiddleware(mux)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "user_service:8080", opts); err != nil {
		log.Error(ctx, "Failed register users grpc service", zap.Error(err))
		return
	}
	log.Info(ctx, "Starting proxy to User")

	if err := pb.RegisterAnalyticHandlerFromEndpoint(ctx, mux, "analytic_service:8081", opts); err != nil {
		log.Error(ctx, "Failed register analytic grpc service", zap.Error(err))
		return
	}
	log.Info(ctx, "Starting proxy to Analytic")

	err := mux.HandlePath("GET", "/ping", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		_, err := w.Write([]byte("pong"))
		if err != nil {
			log.Error(ctx, "Failed to write response", zap.Error(err))
			return
		}
	})
	if err != nil {
		log.Error(ctx, "Failed to register user grpc service", zap.Error(err))
		return
	}

	log.Info(ctx, "Starting gateway server on :5555")
	if err := http.ListenAndServe(":5555", muxWithCORS); err != nil {
		log.Fatal(ctx, "Server exited with error:", zap.Error(err))
	}
}
