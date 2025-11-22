#!/bin/bash
set -e

OPTIMAL=60
ITERS=10

mkdir -p metrics/static_norm

for i in $(seq 1 $ITERS); do
  echo "========================================="
  echo "=== Static Test $i/$ITERS ==="
  echo "========================================="
  
  docker compose -p data_base_project \
    -f docker-compose.yml \
    -f docker-compose.batch.yml \
    down --remove-orphans
  
  rm -f metrics/batch_many_user_test/*
  
  BATCH_SIZE=$OPTIMAL POST_TEST_SLEEP=15 \
  docker compose -p data_base_project \
    -f docker-compose.yml \
    -f docker-compose.batch.yml \
    up --no-recreate batch-runner metrics-collector
  
  ts=$(date +%Y%m%d_%H%M%S)
  [ -f metrics/batch_many_user_test/cpu.csv ] && \
    cp metrics/batch_many_user_test/cpu.csv "metrics/static_norm/cpu_${i}_${ts}.csv"
  [ -f metrics/batch_many_user_test/memory.csv ] && \
    cp metrics/batch_many_user_test/memory.csv "metrics/static_norm/memory_${i}_${ts}.csv"
  [ -f metrics/batch_many_user_test/latency_ms.csv ] && \
    cp metrics/batch_many_user_test/latency_ms.csv "metrics/static_norm/latency_${i}_${ts}.csv"
  
  echo "✅ Saved $i"
  sleep 3
done

echo "✅ ФАЗА 2 ЗАВЕРШЕНА!"
