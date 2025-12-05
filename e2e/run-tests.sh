#!/bin/bash
set -euo pipefail

PROJECT_ROOT="/workspace"
cd "$PROJECT_ROOT"

BACKEND_URL=${BACKEND_URL:-http://backend:8000}
WAIT_MAX_ATTEMPTS=${WAIT_MAX_ATTEMPTS:-120}
POST_TEST_SLEEP=${POST_TEST_SLEEP:-0}
TEST_LABEL=${TEST_LABEL:-e2e}

# record start timestamp for metrics collection (write to temp and publish at end)
mkdir -p /metrics || true
START_TS=$(date +%s)
TMP_TS_FILE="/metrics/timestamps.tmp"
{
  echo "START_TS=${START_TS}"
  echo "TEST_LABEL=${TEST_LABEL}"
} > "${TMP_TS_FILE}"

# Phase toggles
# По умолчанию при локальном запуске гоняем только unit + e2e.
# Интеграционные включаются ТОЛЬКО явно через RUN_INTEGRATION=true.
RUN_UNIT=${RUN_UNIT:-true}
RUN_INTEGRATION=${RUN_INTEGRATION:-false}
RUN_E2E=${RUN_E2E:-true}

if [ "$RUN_UNIT" = "true" ]; then
  echo "Running unit tests ..."
  pushd "$PROJECT_ROOT/backend" >/dev/null
  go test -v ./data_base
  go test -v ./service_logic
  popd >/dev/null
else
  echo "Skipping unit tests"
fi

if [ "$RUN_INTEGRATION" = "true" ] || [ "$RUN_E2E" = "true" ]; then
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
fi

if [ "$RUN_INTEGRATION" = "true" ]; then
  echo "Running integration tests ..."
  pushd "$PROJECT_ROOT/backend" >/dev/null
  go test -v ./integration_tests
  popd >/dev/null
else
  echo "Skipping integration tests"
fi

if [ "$RUN_E2E" = "true" ]; then
  echo "Running e2e tests ..."
  pushd "$PROJECT_ROOT/e2e" >/dev/null
  go mod tidy || true
  # Запускаем только тесты из e2e_test.go, исключая бенчмарки и другие тесты
  # -run "^TestRunAPISuite$" - запускает только TestRunAPISuite из e2e_test.go
  # -bench=^$ - явно отключает бенчмарки
  go test -v -run "^TestRunAPISuite$" -bench=^$
  popd >/dev/null
else
  echo "Skipping e2e tests"
fi

echo "All tests finished."

if [ "$POST_TEST_SLEEP" -gt 0 ]; then
  echo "Sleeping ${POST_TEST_SLEEP}s to allow Prometheus to scrape metrics ..."
  sleep "$POST_TEST_SLEEP"
fi

END_TS=$(date +%s)
echo "END_TS=${END_TS}" >> "${TMP_TS_FILE}"
mv "${TMP_TS_FILE}" /metrics/timestamps.env

