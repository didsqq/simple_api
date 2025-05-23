# Используем официальный образ Go как базовый
FROM golang:1.22-alpine as builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем исходники приложения в рабочую директорию
COPY . .
# Скачиваем все зависимости
RUN go mod download

# Собираем приложение
RUN go build -o main ./cmd/main.go

# Начинаем новую стадию сборки на основе минимального образа
FROM alpine:latest

# Добавляем исполняемый файл из первой стадии в корневую директорию контейнера
COPY --from=builder /app/main /main

# COPY --from=builder /app/.env .env
# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["/main"]
