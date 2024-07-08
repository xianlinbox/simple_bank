#!/bin/sh

set -e
echo "run db migrations"
/app/migrate -path /app/db/migration -database "$DB_SOURCE" up

echo "starting server"
/app/main