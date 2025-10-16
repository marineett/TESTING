FROM ubuntu:22.04

# Установка переменных окружения для избежания интерактивных вопросов
ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=UTC

# Обновление и установка базовых пакетов
RUN apt-get update && apt-get install -y \
    curl \
    ca-certificates \
    git \
    openssh-client \
    unzip \
    wget \
    build-essential \
    cmake \
    pkg-config \
    docker.io \
    docker-compose \
    && rm -rf /var/lib/apt/lists/*

# Установка Go с повторными попытками
RUN for i in 1 2 3; do \
        wget --tries=3 --timeout=60 https://go.dev/dl/go1.24.0.linux-amd64.tar.gz && \
        tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz && \
        rm go1.24.0.linux-amd64.tar.gz && \
        break || sleep 10; \
    done

# Настройка PATH для Go
ENV PATH=$PATH:/usr/local/go/bin
ENV CGO_ENABLED=1

WORKDIR /workspace

# Клонирование репозитория (будет переопределяться при сборке)
RUN git clone https://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@github.com/marineett/TESTING.git . || true

# Скрипт запуска
RUN cat >/usr/local/bin/start-project.sh <<'SH' && chmod +x /usr/local/bin/start-project.sh
#!/bin/bash
set -euo pipefail

cd /workspace

# Обновляем код если контейнер запущен с монтированием
if [ -d .git ]; then
    echo "Updating repository..."
    git pull origin main
fi

docker-compose up --build -d

echo "Checking if backend is ready..."
for i in {1..15}; do
    if docker-compose exec -T backend curl -sf --connect-timeout 1 --max-time 2 http://localhost:8000/api/health >/dev/null; then
        echo "Backend is ready!"
        break
    fi
    echo "Waiting for backend... attempt $i/15"
    sleep 3
done

docker-compose ps

echo "Attaching this test container to app-network..."
docker network connect app-network "$HOSTNAME" 2>/dev/null || true
sleep 2

export PATH=$PATH:/usr/local/go/bin

echo "Running tests..."
cd /workspace/backend
go mod tidy
echo "Unit tests..."
go test -v ./data_base
go test -v ./service_logic

echo "Integration tests..."
go test -v ./integration_tests

echo "Running e2e tests..."
cd /workspace/e2e
go mod tidy
go test -v

echo "All tests completed successfully!"
SH

ENTRYPOINT ["/usr/local/bin/start-project.sh"]