#!/bin/sh

# Run database migrations
echo "Running migrations..."
migrate -path=/migrations-dir -database="$dbDsn" up

# Start the API
echo "Starting API..."
exec /usr/local/bin/api -db-dsn="$dbDsn"
