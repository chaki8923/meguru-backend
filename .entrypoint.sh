#!/bin/sh
set -e

echo "Running migrations..."
migrate -path ./migrations -database "$DATABASE_URL" up

echo "Starting app..."
./main
