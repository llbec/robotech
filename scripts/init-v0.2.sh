#!/bin/sh
set -eu
# Run from the repository root; existing credentials and volumes are preserved.
mkdir -p secrets
chmod 700 secrets
for v02_name in trade-log-token postgres-password trade-log-app-password trade-log-migrator-password; do
  if [ ! -f "secrets/$v02_name" ]; then
    (umask 077; openssl rand -hex 32 > "secrets/$v02_name")
  fi
done
v02_app_password=$(cat secrets/trade-log-app-password)
v02_migration_password=$(cat secrets/trade-log-migrator-password)
# These two generated values are interpolated into SQL: accept only hex, no SQL text.
case "$v02_app_password" in ''|*[!0-9a-fA-F]*) printf '%s\n' 'Invalid app password file' >&2; exit 2;; esac
case "$v02_migration_password" in ''|*[!0-9a-fA-F]*) printf '%s\n' 'Invalid migration password file' >&2; exit 2;; esac
if [ ! -f secrets/trade-log-database-url ]; then
  (umask 077; printf 'postgresql://trade_log_app:%s@postgres:5432/robotech\n' "$v02_app_password" > secrets/trade-log-database-url)
fi
if [ ! -f secrets/trade-log-migration-url ]; then
  (umask 077; printf 'postgresql://trade_log_migrator:%s@postgres:5432/robotech\n' "$v02_migration_password" > secrets/trade-log-migration-url)
fi
chmod 600 secrets/trade-log-app-password secrets/trade-log-migrator-password
chmod 644 secrets/trade-log-token secrets/postgres-password secrets/trade-log-database-url secrets/trade-log-migration-url
docker compose up -d --wait --wait-timeout 90 postgres
# SQL is piped through stdin, passwords are not printed or put in command arguments.
docker compose exec -T postgres psql -q -v ON_ERROR_STOP=1 -U robotech_admin -d robotech <<SQL
SELECT format('CREATE ROLE trade_log_migrator LOGIN PASSWORD %L', '$v02_migration_password') WHERE NOT EXISTS (SELECT FROM pg_roles WHERE rolname='trade_log_migrator')\gexec
SELECT format('CREATE ROLE trade_log_app LOGIN PASSWORD %L', '$v02_app_password') WHERE NOT EXISTS (SELECT FROM pg_roles WHERE rolname='trade_log_app')\gexec
GRANT CONNECT, CREATE ON DATABASE robotech TO trade_log_migrator;
GRANT CREATE, USAGE ON SCHEMA public TO trade_log_migrator;
GRANT CONNECT ON DATABASE robotech TO trade_log_app;
SQL
unset v02_app_password v02_migration_password
docker compose --profile init run --build --rm evidence-init
docker compose --profile init run --rm trade-log-migrate
