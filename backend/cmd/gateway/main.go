package main

import (
	pb "University-Selection-Service/pkg/api"
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	defer logger.Sync()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "user_service:8080", opts); err != nil {
		logger.Error("Failed register user grpc service", zap.Error(err))
		return
	}

	mux.HandlePath("GET", "/ping", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.Write([]byte("pong"))
	})

	logger.Info("Starting gateway server on :5555")
	if err := http.ListenAndServe(":5555", mux); err != nil {
		logger.Fatal("Server exited with error", zap.Error(err))
	}
}
