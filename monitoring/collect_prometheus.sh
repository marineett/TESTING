#!/bin/sh
set -eu

PROM_URL=${PROM_URL:-http://prometheus:9090}
STEP=${STEP:-5s}
OUTPUT_DIR=${OUTPUT_DIR:-/metrics}
# Regex for cadvisor container label 'name'
NAME_REGEX=${NAME_REGEX:-'.*backend.*'}
MARGIN_SECONDS=${MARGIN_SECONDS:-60}

mkdir -p "${OUTPUT_DIR}" || true

if [ -f "${OUTPUT_DIR}/timestamps.env" ]; then
  # shellcheck disable=SC1090
  . "${OUTPUT_DIR}/timestamps.env"
fi

# If END_TS is missing, wait a bit for test runner to finalize
if [ -z "${END_TS:-}" ]; then
  echo "Waiting for END_TS in timestamps.env ..."
  for i in $(seq 1 120); do
    if [ -f "${OUTPUT_DIR}/timestamps.env" ] && grep -q '^END_TS=' "${OUTPUT_DIR}/timestamps.env"; then
      # shellcheck disable=SC1090
      . "${OUTPUT_DIR}/timestamps.env"
      break
    fi
    sleep 2
  done
fi

if [ -z "${START_TS:-}" ] || [ -z "${END_TS:-}" ]; then
  echo "WARN: START_TS/END_TS not found, defaulting to last 10 minutes"
  END_TS=$(date +%s)
  START_TS=$((END_TS - 600))
fi

END_Q=$((END_TS + MARGIN_SECONDS))
echo "Collecting Prometheus metrics from ${START_TS} to ${END_Q} (step ${STEP})"

# Allow overriding full PromQL via env; default to queries you use in GUI
PROMQL_CPU_DEFAULT="sum(rate(container_cpu_usage_seconds_total{job=\"cadvisor\", name=~\"${NAME_REGEX}\"}[120s]))"
PROMQL_MEM_DEFAULT="max_over_time(container_memory_working_set_bytes{job=\"cadvisor\", name=~\"${NAME_REGEX}\"}[5m])"
q_cpu=${PROMQL_CPU:-$PROMQL_CPU_DEFAULT}
q_mem=${PROMQL_MEM:-$PROMQL_MEM_DEFAULT}

fetch_range_to_csv() {
  query="$1"
  outfile="$2"
  curl -sS -X POST "${PROM_URL}/api/v1/query_range" \
    --data-urlencode "query=${query}" \
    --data-urlencode "start=${START_TS}" \
    --data-urlencode "end=${END_Q}" \
    --data-urlencode "step=${STEP}" \
  | jq -r '
      .data.result[] as $s
      | $s.values[]
      | [($s.metric.name // "total"), (.[0]|tostring), (.[1]|tostring)]
      | @csv
    ' > "${outfile}" || {
      echo "WARN: failed to fetch ${outfile}, writing empty file" >&2
      : > "${outfile}"
    }
  echo "Wrote ${outfile}"
}

fetch_range_to_csv "${q_cpu}"  "${OUTPUT_DIR}/cpu.csv"
fetch_range_to_csv "${q_mem}"  "${OUTPUT_DIR}/memory.csv"

echo "Done. Files in ${OUTPUT_DIR}:"
ls -l "${OUTPUT_DIR}" || true


