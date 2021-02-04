#!/bin/bash -e

exec > >(tee -a /var/log/app/entry.log|logger -t server -s 2>/dev/console) 2>&1

echo "[`date`] Running entrypoint script in the '${APP_ENV}' environment..."

CONFIG_FILE=./config/docker.yml

echo "[`date`] Running DB migrations..."
migrate -database "${APP_DSN}" -path ./migrations up

echo "[`date`] Starting server..."
./server -config ${CONFIG_FILE} >> /var/log/app/server.log 2>&1
