#!/usr/bin/env sh
set -e

VAULT_ADDR="${VAULT_ADDR:-http://vault:8200}"
VAULT_TOKEN="${VAULT_TOKEN:?VAULT_TOKEN is required}"
VAULT_SECRET_PATH="${VAULT_SECRET_PATH:-secret/data/postgres}"

echo "Fetching Postgres credentials from Vault at ${VAULT_SECRET_PATH}..."

JSON=$(curl -sS \
  -H "X-Vault-Token: ${VAULT_TOKEN}" \
  "${VAULT_ADDR}/v1/${VAULT_SECRET_PATH}")

export POSTGRES_USER=$(echo "$JSON" | jq -r '.data.data.PG_USERNAME')
export POSTGRES_PASSWORD=$(echo "$JSON" | jq -r '.data.data.PG_PASSWORD')
export POSTGRES_DB=$(echo "$JSON" | jq -r '.data.data.PG_DATABASE')

if [ -z "$POSTGRES_USER" ] || [ -z "$POSTGRES_PASSWORD" ] || [ -z "$POSTGRES_DB" ]; then
  echo "ERROR: missing fields in Vault secret (expected user/password/database)" >&2
  exit 1
fi

echo "Postgres credentials loaded from Vault, starting postgres..."

echo "DEBUG: POSTGRES_USER=$POSTGRES_USER"
echo "DEBUG: POSTGRES_PASSWORD=$POSTGRES_PASSWORD"
echo "DEBUG: POSTGRES_DB=$POSTGRES_DB"

exec docker-entrypoint.sh "$@"
