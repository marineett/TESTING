#!/bin/bash
set -euo pipefail

PROJECT_ROOT="/workspace"
cd "$PROJECT_ROOT"

BACKEND_URL=${BACKEND_URL:-http://backend:8000}
WAIT_MAX_ATTEMPTS=${WAIT_MAX_ATTEMPTS:-120}

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

echo "Running unit tests ..."
pushd "$PROJECT_ROOT/backend" >/dev/null
go test -v ./data_base
go test -v ./service_logic

echo "Running integration tests ..."
go test -v ./integration_tests
popd >/dev/null

echo "Running e2e tests ..."
pushd "$PROJECT_ROOT/e2e" >/dev/null
go mod tidy || true
go test -v
popd >/dev/null

echo "All tests finished."


