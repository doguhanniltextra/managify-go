#!/usr/bin/env bash
set -euo pipefail

TARGET_URL="${1:-http://localhost:8080/health}"
MAX_RETRIES="${2:-12}"
SLEEP_SECONDS="${3:-5}"

echo "Starting health check against: ${TARGET_URL}"
echo "Max retries: ${MAX_RETRIES}, Sleep between retries: ${SLEEP_SECONDS}s"

for i in $(seq 1 "${MAX_RETRIES}"); do
  echo "[Attempt ${i}/${MAX_RETRIES}] Checking service health..."
  STATUS_CODE=$(curl -s -o /dev/null -w "%{http_code}" --max-time 5 "${TARGET_URL}" || echo "000")

  if [ "${STATUS_CODE}" -eq 200 ]; then
    echo "Health check passed! Service is UP and healthy (HTTP 200)."
    exit 0
  fi

  echo "Service returned HTTP ${STATUS_CODE}. Waiting ${SLEEP_SECONDS}s..."
  sleep "${SLEEP_SECONDS}"
done

echo "Health check failed after ${MAX_RETRIES} attempts!"
exit 1
