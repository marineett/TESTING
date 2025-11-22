#!/usr/bin/env sh
set -eu

# Config (override via env if needed)
ITERATIONS="${ITERATIONS:-20}"

# Compose files for degradation pipeline
BASE_COMPOSE="-f docker-compose.yml -f docker-compose.degradation.yml"

# Metrics directories (as used by the runner/collector)
METRICS_DIR="metrics/degradation_find_test"
CPU_DIR="${METRICS_DIR}/cpu"
MEM_DIR="${METRICS_DIR}/memory"
LAT_DIR="${METRICS_DIR}/latency_ms"

mkdir -p "${CPU_DIR}" "${MEM_DIR}" "${LAT_DIR}"

echo "Running degradation tests ${ITERATIONS} times"

i=1
while [ "${i}" -le "${ITERATIONS}" ]; do
  ts="$(date +%Y%m%d_%H%M%S)"
  echo "=== Iteration ${i}/${ITERATIONS} @ ${ts} ==="

  # Ensure clean stack before starting (server from zero)
  docker compose ${BASE_COMPOSE} down --remove-orphans || true

  # Ensure fresh timestamps per iteration
  rm -f "${METRICS_DIR}/timestamps.env" 2>/dev/null || true

  # Run degradation runner and collector with full dependency stack
  docker compose ${BASE_COMPOSE} up --build --force-recreate --abort-on-container-exit degradation-runner metrics-collector || true

  # Copy CSVs into separate folders with unique names
  if [ -f "${METRICS_DIR}/cpu.csv" ]; then
    cp "${METRICS_DIR}/cpu.csv" "${CPU_DIR}/cpu_${ts}_iter${i}.csv"
    echo "Saved CPU -> ${CPU_DIR}/cpu_${ts}_iter${i}.csv"
  else
    echo "WARN: ${METRICS_DIR}/cpu.csv not found"
  fi

  if [ -f "${METRICS_DIR}/memory.csv" ]; then
    cp "${METRICS_DIR}/memory.csv" "${MEM_DIR}/memory_${ts}_iter${i}.csv"
    echo "Saved Memory -> ${MEM_DIR}/memory_${ts}_iter${i}.csv"
  else
    echo "WARN: ${METRICS_DIR}/memory.csv not found"
  fi

  if [ -f "${METRICS_DIR}/latency_ms.csv" ]; then
    cp "${METRICS_DIR}/latency_ms.csv" "${LAT_DIR}/latency_ms_${ts}_iter${i}.csv"
    echo "Saved Latency -> ${LAT_DIR}/latency_ms_${ts}_iter${i}.csv"
  else
    echo "WARN: ${METRICS_DIR}/latency_ms.csv not found"
  fi

  # Clean up containers (keep named volumes)
  docker compose ${BASE_COMPOSE} down --remove-orphans

  i=$((i+1))
done

echo "Done. Results in:"
echo "  ${CPU_DIR}"
echo "  ${MEM_DIR}"
echo "  ${LAT_DIR}"

