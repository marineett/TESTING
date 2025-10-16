FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    curl \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN apt-get update && apt-get install -y git openssh-client unzip wget && rm -rf /var/lib/apt/lists/*

RUN wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz && \
    rm go1.24.0.linux-amd64.tar.gz && \
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /root/.bashrc

# Enable CGO and install build toolchain for go-duckdb
ENV CGO_ENABLED=1
RUN apt-get update && apt-get install -y --no-install-recommends build-essential cmake pkg-config && rm -rf /var/lib/apt/lists/*

RUN apt-get update && apt-get install -y docker.io docker-compose && rm -rf /var/lib/apt/lists/*

WORKDIR /workspace

RUN git clone https://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@github.com/marineett/TESTING.git .

RUN cat >/usr/local/bin/start-project.sh <<'SH' && chmod +x /usr/local/bin/start-project.sh
#!/bin/bash
set -euo pipefail
cd /workspace

# Use unique project name per run to avoid resource clashes
export COMPOSE_PROJECT_NAME=${COMPOSE_PROJECT_NAME:-run_$(date +%s)}

echo "Pruning Docker caches before build..."
# Legacy builder cache prune
docker builder prune -a -f || true
# BuildKit buildx cache prune (if present)
if docker buildx version >/dev/null 2>&1; then
  docker buildx prune -a -f || true
fi
# Remove dangling images/containers/networks/volumes (safe in CI container)
docker system prune -a -f || true

echo "Bringing up stack with fresh, no-cache builds..."
docker-compose down -v --rmi all --remove-orphans || true
docker-compose build --no-cache --pull
docker-compose up -d

echo "Checking if backend is ready (internal check)..."
for i in {1..10}; do
  if docker-compose exec -T backend curl -sf --connect-timeout 1 --max-time 2 http://localhost:8000/api/health >/dev/null; then
    echo "Backend is ready!"
    break
  fi
  echo "Waiting for backend... attempt $i/10"
  sleep 2
done

docker-compose ps

RUNTIME_NET="${COMPOSE_PROJECT_NAME}_app-network"
echo "Attaching this test container to ${RUNTIME_NET}..."
docker network connect "$RUNTIME_NET" "$HOSTNAME" || true
sleep 1

export PATH=$PATH:/usr/local/go/bin

cd /workspace/backend
go mod tidy
echo "unit tests..."
go test -v ./data_base
go test -v ./service_logic

echo "integration tests..."
go test -v ./integration_tests

echo "Running e2e tests..."
cd /workspace/e2e
go mod tidy
go test -v

echo "Tests completed. Press Ctrl+C to stop all services"
# Detach from network before down to avoid 'active endpoints' error
docker network disconnect "$RUNTIME_NET" "$HOSTNAME" || true
docker-compose down -v --remove-orphans || true
SH



ENTRYPOINT ["/usr/local/bin/start-project.sh"]