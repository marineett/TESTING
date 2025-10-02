FROM golang:1.23-alpine

WORKDIR /build

# Копируем файлы go.mod и go.sum
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

EXPOSE 8000

# Используем go run для development
CMD ["go", "run", "main.go"] 