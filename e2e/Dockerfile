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

RUN git clone https://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@github.com/DaniilBelyy2/data_base_project.git .

RUN cat >/usr/local/bin/start-project.sh <<'SH' && chmod +x /usr/local/bin/start-project.sh
#!/bin/bash
set -euo pipefail
cd /workspace
docker-compose up --build -d

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

echo "Attaching this test container to app-network..."
docker network connect app-network "$HOSTNAME" || true
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
wait
SH

ENTRYPOINT ["/usr/local/bin/start-project.sh"]