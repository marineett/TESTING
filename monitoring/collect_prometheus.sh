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
echo "START_TS=$(date -d "@${START_TS}" '+%Y-%m-%d %H:%M:%S' 2>/dev/null || echo "${START_TS}")"
echo "END_TS=$(date -d "@${END_TS}" '+%Y-%m-%d %H:%M:%S' 2>/dev/null || echo "${END_TS}")"
echo "END_Q=$(date -d "@${END_Q}" '+%Y-%m-%d %H:%M:%S' 2>/dev/null || echo "${END_Q}")"
echo "PROM_URL=${PROM_URL}"

# Диагностика: проверим, какие метрики доступны в cAdvisor
echo "=== Diagnostics: Checking available metrics in cAdvisor ==="
echo "Checking if cAdvisor is scraping metrics..."
cadvisor_check=$(curl -sS "${PROM_URL}/api/v1/query?query=up{job=\"cadvisor\"}" | jq -r '.data.result[0].value[1]' 2>/dev/null || echo "failed")
echo "cAdvisor up status: ${cadvisor_check}"
echo ""

echo "Total container_cpu_usage_seconds_total series:"
cpu_count=$(curl -sS "${PROM_URL}/api/v1/query?query=count(container_cpu_usage_seconds_total{job=\"cadvisor\"})" | jq -r '.data.result[0].value[1]' 2>/dev/null || echo "0")
echo "${cpu_count}"
echo ""

echo "Sample container metrics (first 5):"
curl -sS "${PROM_URL}/api/v1/query?query=container_cpu_usage_seconds_total{job=\"cadvisor\"}" | jq -r '.data.result[0:5] | .[] | "name: \(.metric.name // "unknown") | id: \(.metric.id // "unknown") | service: \(.metric.container_label_com_docker_compose_service // "no-label")"' 2>/dev/null || echo "Query failed or no data"
echo ""

echo "All available labels for container metrics:"
curl -sS "${PROM_URL}/api/v1/label/__name__/values" | jq -r '.data[]' | grep -i container | head -10 || echo "No container metrics found"
echo ""

echo "Checking if backend container exists by querying all containers:"
curl -sS "${PROM_URL}/api/v1/query?query=container_cpu_usage_seconds_total{job=\"cadvisor\"}" | jq -r '.data.result[] | select(.metric.name | contains("backend") or .metric.container_label_com_docker_compose_service == "backend") | "Found: name=\(.metric.name // "unknown"), service=\(.metric.container_label_com_docker_compose_service // "no-label"), id=\(.metric.id // "unknown")"' 2>/dev/null | head -5 || echo "No backend containers found"
echo ""

echo "=== End diagnostics ==="
echo ""

# Allow overriding full PromQL via env; default queries
# Используем 60s окно для rate() - баланс между точностью и требованием к длительности теста
# Для памяти используем текущее значение (не max_over_time), чтобы работало на любых тестах
PROMQL_CPU_DEFAULT="sum(rate(container_cpu_usage_seconds_total{job=\"cadvisor\", name=~\"${NAME_REGEX}\"}[60s]))"
PROMQL_MEM_DEFAULT="sum(container_memory_working_set_bytes{job=\"cadvisor\", name=~\"${NAME_REGEX}\"})"
q_cpu=${PROMQL_CPU:-$PROMQL_CPU_DEFAULT}
q_mem=${PROMQL_MEM:-$PROMQL_MEM_DEFAULT}

fetch_range_to_csv() {
  query="$1"
  outfile="$2"
  echo "Fetching ${outfile} with query: ${query}"
  response=$(curl -sS -w "\nHTTP_CODE:%{http_code}" -X POST "${PROM_URL}/api/v1/query_range" \
    --data-urlencode "query=${query}" \
    --data-urlencode "start=${START_TS}" \
    --data-urlencode "end=${END_Q}" \
    --data-urlencode "step=${STEP}") || {
    echo "ERROR: curl failed for ${outfile}" >&2
    : > "${outfile}"
    return 1
  }
  http_code=$(echo "${response}" | grep "HTTP_CODE:" | sed 's/.*HTTP_CODE://')
  body=$(echo "${response}" | sed '/HTTP_CODE:/d')
  if [ "${http_code}" != "200" ]; then
    echo "ERROR: Prometheus returned HTTP ${http_code} for ${outfile}" >&2
    echo "Response: ${body}" >&2
    : > "${outfile}"
    return 1
  fi
  if ! echo "${body}" | jq -e '.data.result' >/dev/null 2>&1; then
    echo "ERROR: Invalid JSON response from Prometheus for ${outfile}" >&2
    echo "Response: ${body}" >&2
    : > "${outfile}"
    return 1
  fi
  result_count=$(echo "${body}" | jq '.data.result | length')
  if [ "${result_count}" = "0" ]; then
    echo "WARN: No data returned from Prometheus for ${outfile} (query returned 0 results)" >&2
    : > "${outfile}"
    return 0
  fi
  echo "${body}" | jq -r '
      .data.result[] as $s
      | $s.values[]
      | [($s.metric.name // "total"), (.[0]|tostring), (.[1]|tostring)]
      | @csv
    ' > "${outfile}" || {
      echo "ERROR: jq processing failed for ${outfile}" >&2
      : > "${outfile}"
      return 1
    }
  line_count=$(wc -l < "${outfile}" | tr -d ' ')
  echo "Wrote ${outfile} (${line_count} lines, ${result_count} series)"
}

fetch_range_to_csv "${q_cpu}"  "${OUTPUT_DIR}/cpu.csv"
fetch_range_to_csv "${q_mem}"  "${OUTPUT_DIR}/memory.csv"

echo "Done. Files in ${OUTPUT_DIR}:"
ls -l "${OUTPUT_DIR}" || true


