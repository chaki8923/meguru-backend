#!/bin/sh
set -e

echo "Running migrations..."
migrate -path ./scripts/db/migrations -database "$DATABASE_URL" up

echo "Starting app..."
exec "$@"
