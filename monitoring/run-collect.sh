#!/usr/bin/env sh
set -eu

# Simple wrapper to run collect_prometheus.sh inside a disposable container.
# Env overrides (optional):
#   NETWORK        - docker network name (default: data_base_project_app-network)
#   OUTPUT_DIR     - output dir inside mounted ./metrics (default: $PWD/metrics/degradation)
#   PROM_URL       - Prometheus URL from the container (default: http://prometheus:9090)
#   NAME_REGEX     - regex to match backend container name label (default: .*backend.*)
#   STEP           - scrape step (default: 5s)
#   PROMQL_CPU     - PromQL for CPU
#   PROMQL_MEM     - PromQL for memory
#
# Usage:
#   sh monitoring/run-collect.sh
#   OUTPUT_DIR="$PWD/metrics/degradation" NAME_REGEX=".*backend.*" sh monitoring/run-collect.sh
#   NETWORK=my_net sh monitoring/run-collect.sh

# Hardcoded settings for this project
NETWORK="data_base_project_app-network"
HOST_OUTPUT_DIR="$PWD/metrics/degradation_many_user_test"
CONTAINER_OUTPUT_DIR="/metrics/degradation_many_user_test"
PROM_URL="http://prometheus:9090"
STEP="5s"
PROMQL_CPU='sum(rate(container_cpu_usage_seconds_total{job="cadvisor",name=~".backend-1."}[120s]))'
PROMQL_MEM='max_over_time(container_memory_working_set_bytes{job="cadvisor",name=~".backend-1."}[5m])'

mkdir -p "$HOST_OUTPUT_DIR"

docker run --rm \
  --network "$NETWORK" \
  -v "$PWD/metrics:/metrics" \
  -v "$PWD/monitoring/collect_prometheus.sh:/collect.sh:ro" \
  alpine:3.19 sh -lc '
    apk add --no-cache curl jq >/dev/null;
    OUTPUT_DIR="'"$CONTAINER_OUTPUT_DIR"'" \
    PROM_URL="'"$PROM_URL"'" \
    STEP="'"$STEP"'" \
    PROMQL_CPU="'"$PROMQL_CPU"'" \
    PROMQL_MEM="'"$PROMQL_MEM"'" \
    sh /collect.sh
  '


