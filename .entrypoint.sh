#!/bin/sh
set -e

echo "Starting Meguru Backend..."
echo "Database URL: $(echo $DATABASE_URL | sed 's/:.*@/:***@/')"

# Run database migrations
echo "Running database migrations..."
if [ -n "$DATABASE_URL" ]; then
    migrate -path ./migrations -database "$DATABASE_URL" up
    echo "Migrations completed successfully"
else
    echo "Warning: DATABASE_URL not set, skipping migrations"
fi

echo "Starting application server..."
exec ./main
