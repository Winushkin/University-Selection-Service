package main

import (
	"University-Selection-Service/internal/analytic/repository"
	analytic "University-Selection-Service/internal/analytic/server"
	"University-Selection-Service/internal/config"
	"University-Selection-Service/internal/interceptors"
	"University-Selection-Service/pkg/api"
	"University-Selection-Service/pkg/logger"
	"University-Selection-Service/pkg/resilence"
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
)

// main starts analytic service
func main() {
	ctx := context.Background()
	ctx, _ = logger.NewLogger(ctx)
	log := logger.GetLoggerFromCtx(ctx)

	cfg, err := config.NewAnalyticCfg()
	if cfg == nil || err != nil {
		log.Error(ctx, "failed to load configuration", zap.Error(err))
		return
	}

	log.Info(ctx, "Successfully connected to university postgres")

	lis, err := net.Listen("tcp", ":"+cfg.RESTPort)
	if err != nil {
		log.Error(ctx, "failed to listen", zap.Error(err))
		return
	}
	log.Info(ctx, "Listening on", zap.String("port", cfg.RESTPort))

	conn, err := grpc.NewClient("user_service:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(ctx, "failed to connect to user server", zap.Error(err))
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Error(ctx, "failed to close client connection", zap.Error(err))
		}
	}(conn)
	client := api.NewUserServiceClient(conn)

	r, err := repository.NewAnalyticRepository(ctx, cfg.Postgres)
	if err != nil {
		log.Error(ctx, "failed to create repository", zap.Error(err))
		return
	}

	srv, err := analytic.New(client, r)
	if err != nil {
		log.Error(ctx, "failed to create analytic service", zap.Error(err))
		return
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.AuthInterceptor(cfg.JWTSecret)),
	)

	api.RegisterAnalyticServer(server, srv)
	reflection.Register(server)

	if err = resilence.Retry(func() error { return server.Serve(lis) }, 5, 100); err != nil {
		log.Error(ctx, "failed to serve analytic server", zap.Error(err))
		return
	}

	healthServer := health.NewServer()

	grpc_health_v1.RegisterHealthServer(server, healthServer)

	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
}
