FROM golang:1.24 AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/cmd/user/ ./backend/cmd/user/
COPY backend/internal/ ./backend/internal/
COPY backend/pkg/ ./backend/pkg/
COPY backend/db/ ./backend/db/

RUN go build -o user ./backend/cmd/user/main.go
ADD https://github.com/grpc-ecosystem/grpc-health-probe/releases/latest/download/grpc_health_probe-linux-amd64 /bin/grpc_health_probe
RUN chmod +x /bin/grpc_health_probe

FROM debian:bookworm-slim
WORKDIR /app

COPY --from=builder /app/user ./

ENTRYPOINT ["./user"]