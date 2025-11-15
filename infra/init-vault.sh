#!/usr/bin/env bash

# TODO Прописать в ридми запуск
set -e

echo "Starting Vault..."
docker compose up -d vault

echo "Waiting for Vault to be ready..."
sleep 5

export VAULT_ADDR=http://localhost:8201
export VAULT_TOKEN=root

echo "Initializing secrets in Vault..."
vault kv put secret/avito-test-task \
  PG_USERNAME=postgres \
  PG_PASSWORD=postgres \
  PG_DATABASE=avito-test-task-db \
  PG_HOST=postgres \
  PG_PORT=5432 \
  APP_PORT=8080

echo "Secrets initialized. You can now run: docker compose up"
