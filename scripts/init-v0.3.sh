#!/bin/sh
set -eu
# Run from the repository root; existing credentials and volumes are preserved.
mkdir -p secrets
chmod 700 secrets
for v03_name in trade-log-token postgres-password trade-log-app-password trade-log-migrator-password trade-log-collector-password; do
  if [ ! -f "secrets/$v03_name" ]; then
    (umask 077; openssl rand -hex 32 > "secrets/$v03_name")
  fi
done
v03_app_password=$(cat secrets/trade-log-app-password)
v03_migration_password=$(cat secrets/trade-log-migrator-password)
v03_collector_password=$(cat secrets/trade-log-collector-password)
# These generated values are interpolated into SQL: accept only hex, no SQL text.
case "$v03_app_password" in ''|*[!0-9a-fA-F]*) printf '%s\n' 'Invalid app password file' >&2; exit 2;; esac
case "$v03_migration_password" in ''|*[!0-9a-fA-F]*) printf '%s\n' 'Invalid migration password file' >&2; exit 2;; esac
case "$v03_collector_password" in ''|*[!0-9a-fA-F]*) printf '%s\n' 'Invalid collector password file' >&2; exit 2;; esac
if [ ! -f secrets/trade-log-collector-database-url ]; then
  (umask 077; printf 'postgresql://trade_log_collector:%s@postgres:5432/robotech\n' "$v03_collector_password" > secrets/trade-log-collector-database-url)
fi
if [ ! -f secrets/trade-log-database-url ]; then
  (umask 077; printf 'postgresql://trade_log_app:%s@postgres:5432/robotech\n' "$v03_app_password" > secrets/trade-log-database-url)
fi
if [ ! -f secrets/trade-log-migration-url ]; then
  (umask 077; printf 'postgresql://trade_log_migrator:%s@postgres:5432/robotech\n' "$v03_migration_password" > secrets/trade-log-migration-url)
fi
chmod 600 secrets/trade-log-app-password secrets/trade-log-migrator-password secrets/trade-log-collector-password
chmod 644 secrets/trade-log-token secrets/postgres-password secrets/trade-log-database-url secrets/trade-log-migration-url secrets/trade-log-collector-database-url
docker compose up -d --wait --wait-timeout 90 postgres
# SQL is piped through stdin, passwords are not printed or put in command arguments.
docker compose exec -T postgres psql -q -v ON_ERROR_STOP=1 -U robotech_admin -d robotech <<SQL
SELECT format('CREATE ROLE trade_log_migrator LOGIN PASSWORD %L', '$v03_migration_password') WHERE NOT EXISTS (SELECT FROM pg_roles WHERE rolname='trade_log_migrator')\gexec
SELECT format('CREATE ROLE trade_log_app LOGIN PASSWORD %L', '$v03_app_password') WHERE NOT EXISTS (SELECT FROM pg_roles WHERE rolname='trade_log_app')\gexec
SELECT format('CREATE ROLE trade_log_collector LOGIN PASSWORD %L', '$v03_collector_password') WHERE NOT EXISTS (SELECT FROM pg_roles WHERE rolname='trade_log_collector')\gexec
GRANT CONNECT ON DATABASE robotech TO trade_log_collector;
GRANT CONNECT, CREATE ON DATABASE robotech TO trade_log_migrator;
GRANT CREATE, USAGE ON SCHEMA public TO trade_log_migrator;
GRANT CONNECT ON DATABASE robotech TO trade_log_app;
SQL
unset v03_app_password v03_migration_password v03_collector_password
docker compose --profile init run --build --rm evidence-init
docker compose --profile init run --rm trade-log-migrate
