#!/bin/bash
set -euo pipefail

PROJECT_ROOT=${PROJECT_ROOT:-/workspace}

set +e
/usr/local/bin/run-tests.sh
STATUS=$?
set -e

export PROJECT_ROOT
export CLEAN=${CLEAN:-false}
/usr/local/bin/collect-allure.sh

echo "Tests finished and Allure results collected to $PROJECT_ROOT/e2e/allure-results"

exit ${STATUS}


