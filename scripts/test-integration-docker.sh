#!/usr/bin/env bash
set -euo pipefail

# Run Postgres-backed integration tests against a temporary Docker container,
# even if local 5432 is occupied (uses random host port via `-P`).
# Usage:
#   bash scripts/test-integration-docker.sh           # default: go test -v ./test -run Postgres
#   bash scripts/test-integration-docker.sh ./...     # custom go test target

TEST_TARGET="${1:-./test}"
shift || true
RUN_ARGS=("$@")

echo "[info] Starting temporary Postgres container for integration tests..."
CID=$(docker run -d \
  -e POSTGRES_USER=app \
  -e POSTGRES_PASSWORD=app_password \
  -e POSTGRES_DB=productdb \
  -P \
  postgres:16-alpine)

cleanup() {
  echo "[info] Cleaning up container $CID" >&2
  docker rm -f "$CID" >/dev/null 2>&1 || true
}
trap cleanup EXIT

echo "[info] Waiting for Postgres to be ready..."
for i in {1..60}; do
  if docker exec "$CID" pg_isready -h 127.0.0.1 -p 5432 -U app -d productdb >/dev/null 2>&1; then
    echo "[ok] Postgres ready"
    break
  fi
  sleep 1
done

if ! docker exec "$CID" pg_isready -h 127.0.0.1 -p 5432 -U app -d productdb >/dev/null 2>&1; then
  echo "[error] Postgres did not become ready in time" >&2
  exit 1
fi

# Apply migrations from repo using migrate container within the same netns
MIGRATIONS_DIR="$(pwd)/apps/product-query-svc/adapters/outbound/postgres/migrations"
echo "[info] Applying migrations from: ${MIGRATIONS_DIR}"
docker run --rm \
  --network container:"$CID" \
  -v "${MIGRATIONS_DIR}:/migrations:ro" \
  migrate/migrate \
  -path /migrations -database "postgres://app:app_password@127.0.0.1:5432/productdb?sslmode=disable" up

# Discover host port for container's 5432
HOST_PORT=$(docker port "$CID" 5432/tcp | sed -E 's/.*:(\d+)/\1/' | head -n1)
if [[ -z "${HOST_PORT}" ]]; then
  echo "[error] Failed to resolve host port for Postgres" >&2
  exit 1
fi
export DATABASE_URL="postgres://app:app_password@127.0.0.1:${HOST_PORT}/productdb?sslmode=disable"
echo "[info] DATABASE_URL=${DATABASE_URL}"

echo "[info] Running tests: go test -v ${TEST_TARGET} ${RUN_ARGS[*]:-}"
go test -v "${TEST_TARGET}" "${RUN_ARGS[@]:-}"

