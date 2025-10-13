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

EXPOSE 8000
CMD ["go", "run", "main.go"]