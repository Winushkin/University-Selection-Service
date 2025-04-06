package main

import (
	"University-Selection-Service/internal/server"
	logger2 "University-Selection-Service/pkg/logger"
	"context"
)

func main() {
	ctx := context.Background()
	ctx, _ = logger2.NewLogger(ctx)
	server.StartServer(ctx)
	//test commit 6
}
