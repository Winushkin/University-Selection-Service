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

func main() {
	ctx := context.Background()
	ctx, _ = logger.NewLogger(ctx)
	log := logger.GetLoggerFromCtx(ctx)

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "user_service:8080", opts); err != nil {
		log.Error(ctx, "Failed register user grpc service", zap.Error(err))
		return
	}

	mux.HandlePath("GET", "/ping", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.Write([]byte("pong"))
	})

	log.Info(ctx, "Starting gateway server on :5555")
	if err := http.ListenAndServe(":5555", mux); err != nil {
		log.Fatal(ctx, "Server exited with error:", zap.Error(err))
	}
}
