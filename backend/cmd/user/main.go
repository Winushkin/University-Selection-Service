package main

import (
	"University-Selection-Service/internal/config"
	"University-Selection-Service/internal/interceptors"
	"University-Selection-Service/internal/user"
	"University-Selection-Service/pkg/api"
	"University-Selection-Service/pkg/logger"
	"University-Selection-Service/pkg/postgres"
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
)

func main() {
	ctx := context.Background()
	ctx, _ = logger.NewLogger(ctx)
	log := logger.GetLoggerFromCtx(ctx)

	cfg, err := config.New()
	if cfg == nil || err != nil {
		log.Error(ctx, "failed to load configuration", zap.Error(err))
		return
	}

	db, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		log.Error(ctx, "failed to connect to users postgres", zap.Error(err))
		return
	}
	defer db.Close()
	log.Info(ctx, "Successfully connected to users postgres")

	lis, err := net.Listen("tcp", ":"+cfg.RESTPort)
	if err != nil {
		log.Error(ctx, "failed to listen", zap.Error(err))
		return
	}
	log.Info(ctx, "Listening on", zap.String("port", cfg.RESTPort))

	srv, err := user.New(ctx, cfg, cfg.JWTSecret)
	if err != nil {
		log.Error(ctx, "failed to create user service", zap.Error(err))
		return
	}
	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.AuthInterceptor(cfg)),
	)
	api.RegisterUserServiceServer(server, srv)
	reflection.Register(server)
	if err := server.Serve(lis); err != nil {
		log.Error(ctx, "failed to serve gRPC server", zap.Error(err))
	}
	healthServer := health.NewServer()

	grpc_health_v1.RegisterHealthServer(server, healthServer)

	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

}
