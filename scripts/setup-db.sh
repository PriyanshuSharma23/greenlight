#!/bin/bash
set -e

# Wait until PostgreSQL is ready
until pg_isready -h localhost -p 5432 -U "$POSTGRES_USER"; do
  echo "Waiting for PostgreSQL to start..."
  sleep 1
done


# Create the greenlight database
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE greenlight;
    \c greenlight;
    CREATE EXTENSION citext;
EOSQL

# Run migrations
migrate -path ./migrations -database "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@localhost/greenlight?sslmode=disable" up
