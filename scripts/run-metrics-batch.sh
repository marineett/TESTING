#!/usr/bin/env bash
set -euo pipefail

# Batch runner: starts the stack (base + monitoring), waits for metrics, and stores CPU/MEM CSVs per iteration.
# Usage: scripts/run-metrics-batch.sh [COUNT] [OUTPUT_DIR]
#   COUNT      - number of iterations (default: 100)
#   OUTPUT_DIR - directory to store results (default: metrics/batch-YYYYmmddHHMMSS)

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
COMPOSE_BASE="${PROJECT_ROOT}/docker-compose.yml"
COMPOSE_MON="${PROJECT_ROOT}/docker-compose.monitoring.yml"
METRICS_DIR="${PROJECT_ROOT}/metrics"

COUNT="${1:-100}"
STAMP="$(date +%Y%m%d%H%M%S)"
OUT_DIR="${2:-${METRICS_DIR}/batch-${STAMP}}"

mkdir -p "${OUT_DIR}/cpu" "${OUT_DIR}/memory"

echo "Batch run: COUNT=${COUNT} OUT_DIR=${OUT_DIR}"

require() {
  command -v "$1" >/dev/null 2>&1 || { echo "ERROR: '$1' not found in PATH" >&2; exit 1; }
}

require docker

cleanup_stack() {
  # Keep named volumes (Prometheus TSDB) for speed; remove containers and networks only
  docker compose -f "${COMPOSE_BASE}" -f "${COMPOSE_MON}" down --remove-orphans || true
}

wait_for_files() {
  # Args: since_epoch [timeout_sec]
  local since_epoch="${1:-0}"
  local timeout="${2:-300}"
  local start_ts
  start_ts=$(date +%s)
  while true; do
    if [ -f "${METRICS_DIR}/cpu.csv" ] && [ -f "${METRICS_DIR}/memory.csv" ]; then
      # Ensure files are updated after this iteration started (avoid reusing from previous run)
      local cpu_mtime mem_mtime
      cpu_mtime=$(stat -c '%Y' "${METRICS_DIR}/cpu.csv" 2>/dev/null || echo 0)
      mem_mtime=$(stat -c '%Y' "${METRICS_DIR}/memory.csv" 2>/dev/null || echo 0)
      if [ "${cpu_mtime}" -ge "${since_epoch}" ] && [ "${mem_mtime}" -ge "${since_epoch}" ]; then
        if [ -s "${METRICS_DIR}/cpu.csv" ] && [ -s "${METRICS_DIR}/memory.csv" ]; then
          return 0
        fi
      fi
    fi
    if [ $(( $(date +%s) - start_ts )) -ge "${timeout}" ]; then
      echo "Timeout waiting for fresh metrics CSVs (cpu.csv, memory.csv)" >&2
      return 1
    fi
    sleep 2
  done
}

for i in $(seq 1 "${COUNT}"); do
  printf "\n=== Iteration %d/%d ===\n" "${i}" "${COUNT}"

  # Mark iteration start time to detect freshly written files without deleting previous ones
  ITER_START="$(date +%s)"

  # Start stack
  docker compose -f "${COMPOSE_BASE}" -f "${COMPOSE_MON}" up -d --remove-orphans

  # Optionally ensure Prometheus is up (metrics-collector also waits internally)
  for _ in $(seq 1 60); do
    if curl -sf --max-time 2 http://localhost:9090/-/ready >/dev/null 2>&1; then
      break
    fi
    sleep 1
  done

  # Wait for freshly produced metrics CSVs (mtime >= ITER_START)
  if ! wait_for_files "${ITER_START}" 600; then
    echo "Iteration ${i}: metrics not produced within timeout; capturing logs" >&2
    docker compose -f "${COMPOSE_BASE}" -f "${COMPOSE_MON}" logs --no-color --tail=200 prometheus metrics-collector bench-runner e2e-runner || true
    cleanup_stack
    exit 1
  fi

  # Store results with unique names
  ts="$(date +%Y%m%d-%H%M%S)"
  cp "${METRICS_DIR}/cpu.csv"    "${OUT_DIR}/cpu/${ts}_${i}.csv"
  cp "${METRICS_DIR}/memory.csv" "${OUT_DIR}/memory/${ts}_${i}.csv"
  echo "Saved: ${OUT_DIR}/cpu/${ts}_${i}.csv and memory/${ts}_${i}.csv"

  # Stop stack to reset runners between iterations
  cleanup_stack
done

echo "\nDone. Results saved in: ${OUT_DIR}"


