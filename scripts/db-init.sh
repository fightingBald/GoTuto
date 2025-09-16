#!/usr/bin/env bash
set -euo pipefail

# Simple migration runner that initializes the DB schema and seeds test data
# Usage:
#   DATABASE_URL="postgres://app:app_password@localhost:5432/productdb?sslmode=disable" \
#   bash scripts/db-init.sh
#
# If the `migrate` CLI is installed locally, it will be used.
# Otherwise, this script falls back to `docker run migrate/migrate`.

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
MIGRATIONS_REL="apps/product-query-svc/adapters/postgres/migrations"
MIGRATIONS_DIR="${ROOT_DIR}/${MIGRATIONS_REL}"

DATABASE_URL="${DATABASE_URL:-}"
if [[ -z "${DATABASE_URL}" ]]; then
  DATABASE_URL="postgres://app:app_password@localhost:5432/productdb?sslmode=disable"
  echo "[info] DATABASE_URL not set; using default: ${DATABASE_URL}" >&2
fi

echo "[info] Running migrations from: ${MIGRATIONS_DIR}" >&2

if command -v migrate >/dev/null 2>&1; then
  echo "[info] Using local migrate CLI" >&2
  migrate -path "${MIGRATIONS_DIR}" -database "${DATABASE_URL}" up
else
  echo "[info] Using dockerized migrate (migrate/migrate)" >&2
  docker run --rm \
    -v "${MIGRATIONS_DIR}:/migrations:ro" \
    migrate/migrate \
    -path /migrations -database "${DATABASE_URL}" up
fi

echo "[ok] Database initialized with schema and seed data"

