#!/bin/bash
set -euo pipefail

PROJECT_ROOT="/workspace"
cd "$PROJECT_ROOT"

BACKEND_URL=${BACKEND_URL:-http://backend:8000}
WAIT_MAX_ATTEMPTS=${WAIT_MAX_ATTEMPTS:-120}
E2E_TEST_REGEX=${E2E_TEST_REGEX:-TestRunAPISuite}
POST_TEST_SLEEP=${POST_TEST_SLEEP:-0}
TEST_LABEL=${TEST_LABEL:-bench}
SEQUENTIAL_BENCH=${SEQUENTIAL_BENCH:-false}
TEST1_LABEL=${TEST1_LABEL:-create_batch}
TEST1_REGEX=${TEST1_REGEX:-}
TEST2_LABEL=${TEST2_LABEL:-one_big}
TEST2_REGEX=${TEST2_REGEX:-}

# record start timestamp for metrics collection (write to temp and publish at end)
mkdir -p /metrics || true
START_TS=$(date +%s)
TMP_TS_FILE="/metrics/timestamps.tmp"
{
  echo "START_TS=${START_TS}"
  echo "TEST_LABEL=${TEST_LABEL}"
} > "${TMP_TS_FILE}"

echo "Waiting for backend readiness at ${BACKEND_URL}/api/health ..."
ready=0
for ((i=1; i<=WAIT_MAX_ATTEMPTS; i++)); do
  if curl -sf --connect-timeout 1 --max-time 2 "${BACKEND_URL}/api/health" >/dev/null; then
    echo "Backend is ready"
    ready=1
    break
  fi
  echo "... still waiting ($i/${WAIT_MAX_ATTEMPTS})"
  sleep 2
done

if [ "$ready" -ne 1 ]; then
  echo "Backend did not become ready in time at ${BACKEND_URL}. Exiting." >&2
  exit 1
fi

run_one() {
  lbl="$1"; regex="$2"
  echo "Running e2e tests for ${lbl} with regex: ${regex} ..."
  subdir="/metrics/${lbl}"
  mkdir -p "${subdir}" || true
  START_TS=$(date +%s)
  TMP_TS_FILE="${subdir}/timestamps.tmp"
  {
    echo "START_TS=${START_TS}"
    echo "TEST_LABEL=${lbl}"
  } > "${TMP_TS_FILE}"
  pushd "$PROJECT_ROOT/e2e" >/dev/null
  export METRICS_DIR="${subdir}"
  export LATENCY_FILE="${subdir}/latency_ms.csv"
  go test -v -count=1 -run "${regex}"
  popd >/dev/null
  echo "Tests for ${lbl} finished."
  if [ "$POST_TEST_SLEEP" -gt 0 ]; then
    echo "Sleeping ${POST_TEST_SLEEP}s to allow Prometheus scrape ..."
    sleep "$POST_TEST_SLEEP"
  fi
  END_TS=$(date +%s)
  echo "END_TS=${END_TS}" >> "${TMP_TS_FILE}"
  mv "${TMP_TS_FILE}" "${subdir}/timestamps.env"
  echo "Collecting metrics for ${lbl} ..."
  OUTPUT_DIR="${subdir}" /usr/local/bin/collect_prometheus.sh || true
}

if [ "${SEQUENTIAL_BENCH}" = "true" ]; then
  if [ -z "${TEST1_REGEX}" ] || [ -z "${TEST2_REGEX}" ]; then
    echo "SEQUENTIAL_BENCH=true requires TEST1_REGEX and TEST2_REGEX" >&2
    exit 1
  fi
  run_one "${TEST1_LABEL}" "${TEST1_REGEX}"
  run_one "${TEST2_LABEL}" "${TEST2_REGEX}"
else
  echo "Running e2e tests with regex: ${E2E_TEST_REGEX} ..."
  pushd "$PROJECT_ROOT/e2e" >/dev/null
  go test -v -count=1 -run "${E2E_TEST_REGEX}"
  popd >/dev/null
  echo "Tests finished."
  if [ "$POST_TEST_SLEEP" -gt 0 ]; then
    echo "Sleeping ${POST_TEST_SLEEP}s ..."
    sleep "$POST_TEST_SLEEP"
  fi
  END_TS=$(date +%s)
  echo "END_TS=${END_TS}" >> "${TMP_TS_FILE}"
  mv "${TMP_TS_FILE}" /metrics/timestamps.env
  echo "Collecting metrics ..."
  OUTPUT_DIR="/metrics" /usr/local/bin/collect_prometheus.sh || true
fi


