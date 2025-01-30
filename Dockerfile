FROM golang:1.23.2 AS builder

# Устанавливаем необходимые зависимости для кросс-компиляции
RUN apt-get update && apt-get install -y mingw-w64

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем исходный код в контейнер
COPY . .

# Настраиваем переменные окружения для кросс-компиляции
ENV GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc

# Собираем приложение
RUN go build -o output.exe cmd/main.go