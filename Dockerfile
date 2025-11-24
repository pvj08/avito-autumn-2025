FROM golang:1.23-alpine AS builder

WORKDIR /app

# зависимости
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Сборка бинаря
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/app ./

# Кладём миграции (папка migrations в корне проекта)
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

# Стартуем
CMD ["./app"]
