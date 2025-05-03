FROM golang:1.24-alpine3.21 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o gateway ./cmd/gateway/main.go


FROM alpine:3.21

WORKDIR /app

COPY --from=build /app/gateway ./

ENTRYPOINT ["./gateway"]

