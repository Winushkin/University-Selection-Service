FROM golang:1.24 AS builder

WORKDIR /app

COPY ../ .

RUN go build -o analytic ./cmd/analytic/main.go
ADD https://github.com/grpc-ecosystem/grpc-health-probe/releases/latest/download/grpc_health_probe-linux-amd64 /bin/grpc_health_probe
RUN chmod +x /bin/grpc_health_probe

FROM debian:bookworm-slim
WORKDIR /app

COPY --from=builder /app/analytic ./

ENTRYPOINT ["./analytic"]