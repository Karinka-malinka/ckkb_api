# Стадия сборки
FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем бинарник как /build/ckkbapi
RUN go build -o /build/ckkbapi ./cmd

# Финальный минимальный образ
FROM alpine:3.19

RUN apk add --no-cache ca-certificates
RUN apk add --no-cache tzdata

ENV TZ=Europe/Moscow

WORKDIR /app

# Копируем бинарник из builder и переименовываем в /app/app
COPY --from=builder /build/ckkbapi /app/app

EXPOSE 8005

#COPY internal/lib/config/config.toml /app/config.toml

# Запускаем бинарник по абсолютному пути
CMD ["/app/app"]