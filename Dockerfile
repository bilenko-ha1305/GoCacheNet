# Используйте образ golang для сборки вашего приложения
FROM golang:1.20.2 AS builder

WORKDIR /app
COPY . .
RUN go build -o redis-server

# Второй этап: использование легковесного образа для запуска приложения
FROM alpine:latest

# Установка redis-cli внутри контейнера
RUN apk --update add redis

# Копирование бинарного файла из первого этапа
COPY --from=builder /app/redis-server /usr/local/bin/

# Прослушивание порта 8080
EXPOSE 8080

# Запуск вашего приложения
CMD ["redis-server"]
