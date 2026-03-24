#!/usr/bin/env bash
# Run all migrations. Loads DATABASE_URL from services/api/.env if not set.
set -e
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$ROOT/infra/db"

if [ -z "$DATABASE_URL" ] && [ -f "$ROOT/services/api/.env" ]; then
  export DATABASE_URL=$(grep -E '^DATABASE_URL=' "$ROOT/services/api/.env" | sed 's/^DATABASE_URL=//' | sed 's/^["'\'']//;s/["'\'']$//')
fi

if [ -z "$DATABASE_URL" ]; then
  echo "DATABASE_URL is not set. Set it or add it to services/api/.env"
  exit 1
fi

echo "Running migrations..."
psql "$DATABASE_URL" -f migrations/001_schema.sql
psql "$DATABASE_URL" -f migrations/002_coding_logs.sql
psql "$DATABASE_URL" -f migrations/003_tasks_custom.sql
psql "$DATABASE_URL" -f migrations/004_ideas_linkedin.sql
psql "$DATABASE_URL" -f migrations/005_finances_spend_type.sql
psql "$DATABASE_URL" -f migrations/006_finances_insta_paid_collab_type.sql
echo "All migrations done."
