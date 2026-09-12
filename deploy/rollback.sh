#!/usr/bin/env bash
set -euo pipefail

PREVIOUS_TAG="${1:-}"

if [ -z "${PREVIOUS_TAG}" ]; then
  echo "Error: Previous version tag or commit must be specified for rollback."
  echo "Usage: ./rollback.sh <previous_tag_or_sha>"
  exit 1
fi

echo "Initiating rollback to version: ${PREVIOUS_TAG}..."

# Example: If deploying via docker-compose, Kubernetes, or Render hook
# Here we provide the standard rollback trigger pattern:
if [ -n "${RENDER_ROLLBACK_HOOK_URL:-}" ]; then
  echo "Triggering Render rollback hook..."
  curl -s -X POST "${RENDER_ROLLBACK_HOOK_URL}"
fi

echo "Rollback command executed for target: ${PREVIOUS_TAG}"
