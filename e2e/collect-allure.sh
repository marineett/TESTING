#!/bin/sh
set -eu

PROJECT_ROOT=${PROJECT_ROOT:-/workspace}
TARGET_DIR="$PROJECT_ROOT/e2e/allure-results"

CLEAN=${CLEAN:-false}

if [ "$CLEAN" = "true" ]; then
    rm -rf "$TARGET_DIR"
fi

mkdir -p "$TARGET_DIR"

ALLURE_DIRS=$(find "$PROJECT_ROOT" -type d -name 'allure-results' 2>/dev/null | grep -v "^$TARGET_DIR\(/\|$\)") || true

if [ -z "$ALLURE_DIRS" ]; then
    echo "No allure-results directories found."
    exit 0
fi

COPIED=0
for d in $ALLURE_DIRS; do
    rel=${d#${PROJECT_ROOT}/}
    parent_rel=${rel%/allure-results}
    dest_dir="$TARGET_DIR/$parent_rel"

        echo "Collecting from: $d -> $dest_dir"
    mkdir -p "$dest_dir"

    for ext in json xml txt png svg ; do
        set +e
        FILES=$(find "$d" -maxdepth 1 -type f -name "*.$ext" 2>/dev/null)
        set -e
        if [ -n "$FILES" ]; then
            cp -n $FILES "$dest_dir" 2>/dev/null || true
            COUNT=$(echo "$FILES" | wc -w | tr -d ' ')
            COPIED=$((COPIED + COUNT))
        fi
    done
done

echo "Done"


