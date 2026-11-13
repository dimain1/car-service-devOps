
# Этап сборки
FROM golang:1.27.4-alpine AS builder

WORKDIR /app

# Сначала копируем зависимости для кеширования слоёв
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем исполняемый файл
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o car-service .

# Этап запуска
FROM alpine:3.23

WORKDIR /app

# Создаём непривилегированного пользователя
RUN adduser -D -H appuser

COPY --from=builder /app/car-service .

USER appuser

# HTTP и gRPC
EXPOSE 8080 9090

ENTRYPOINT ["./car-service"]
