# backend/Dockerfile
# syntax=docker/dockerfile:1

## Builder stage
FROM golang:1.21 AS builder

WORKDIR /app

# Copy go mod and sum files
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy source files
COPY backend ./

# Build
RUN go build -o /out/server ./cmd/server.go

## Deploy stage
FROM alpine:latest  AS deploy

WORKDIR /app

COPY --from=builder /out/server /app/server
COPY .env.example /app/ # 環境変数をコピー

EXPOSE 8080

CMD ["./server"]
