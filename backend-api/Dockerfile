# ベースイメージとして最新のGolangイメージを使用
FROM golang:1.21-alpine AS base

# 必要なパッケージのインストール
RUN apk add --no-cache gcc musl-dev git

# 作業ディレクトリの設定
WORKDIR /app

# 依存関係のコピーとダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .

# テストステージ
FROM base AS test
RUN go install github.com/golang/mock/mockgen@v1.6.0
ENV CGO_ENABLED=1
ENV GO_ENV=test
CMD ["go", "test", "-v", "./..."]

# ビルドステージ
FROM base AS builder
RUN go build -o main .

# 本番ステージ
FROM alpine:latest AS production
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./main"]