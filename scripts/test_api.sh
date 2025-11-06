#!/usr/bin/env bash
set -euo pipefail

BASE_URL=${BASE_URL:-http://localhost:8080}

echo "Testing API at $BASE_URL"

# GET /api/config
echo -n "GET /api/config -> "
if curl -sS "$BASE_URL/api/config" | jq . >/dev/null 2>&1; then
  echo "OK"
else
  echo "FAIL"
  exit 1
fi

# GET /api/servers
echo -n "GET /api/servers -> "
if curl -sS "$BASE_URL/api/servers" | jq . >/dev/null 2>&1; then
  echo "OK"
else
  echo "FAIL"
  exit 1
fi

# POST /api/connect (use first server if available)
SERVER_ID=$(curl -sS "$BASE_URL/api/servers" | jq -r '.[0].id // .[0].uuid // empty') || true
if [ -n "$SERVER_ID" ]; then
  echo -n "POST /api/connect -> "
  if curl -sS -X POST -H "Content-Type: application/json" -d "{\"server_id\": \"$SERVER_ID\"}" "$BASE_URL/api/connect" | jq . >/dev/null 2>&1; then
    echo "OK"
  else
    echo "FAIL"
    exit 1
  fi

  echo -n "GET /api/status -> "
  if curl -sS "$BASE_URL/api/status" | jq . >/dev/null 2>&1; then
    echo "OK"
  else
    echo "FAIL"
    exit 1
  fi

  echo -n "GET /api/stats -> "
  if curl -sS "$BASE_URL/api/stats" | jq . >/dev/null 2>&1; then
    echo "OK"
  else
    echo "FAIL"
    exit 1
  fi

  echo -n "POST /api/disconnect -> "
  if curl -sS -X POST "$BASE_URL/api/disconnect" | jq . >/dev/null 2>&1; then
    echo "OK"
  else
    echo "FAIL"
    exit 1
  fi
else
  echo "No servers found in /api/servers; skipping connect/disconnect tests."
fi

echo "All smoke tests passed."
