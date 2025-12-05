FROM golang:1.24-bookworm

ENV CGO_ENABLED=1
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential cmake pkg-config && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# Если в build-контексте присутствует файл .env.docker,
# скопируем его в .env внутри образа, чтобы backend мог
# читать переменные через godotenv.Load().
RUN if [ -f .env.docker ]; then cp .env.docker .env; fi

EXPOSE 8000
CMD ["go", "run", "main.go"]