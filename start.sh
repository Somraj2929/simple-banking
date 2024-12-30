#!/bin/sh

set -e 

echo "run db migrations"
source app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start app"
/app/main
