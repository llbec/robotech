#!/bin/sh
set -eu
# Run from the repository root; existing credentials and volumes are preserved.
mkdir -p secrets
chmod 700 secrets
for v05_name in trade-log-token postgres-password trade-log-app-password trade-log-migrator-password trade-log-collector-password trade-log-publisher-password webhook-token; do
  if [ ! -f "secrets/$v05_name" ]; then
    (umask 077; openssl rand -hex 32 > "secrets/$v05_name")
  fi
done
v05_app_password=$(cat secrets/trade-log-app-password)
v05_migration_password=$(cat secrets/trade-log-migrator-password)
v05_collector_password=$(cat secrets/trade-log-collector-password)
v05_publisher_password=$(cat secrets/trade-log-publisher-password)
# These generated values are interpolated into SQL: accept only hex, no SQL text.
case "$v05_app_password" in ''|*[!0-9a-fA-F]*) printf '%s\n' 'Invalid app password file' >&2; exit 2;; esac
case "$v05_migration_password" in ''|*[!0-9a-fA-F]*) printf '%s\n' 'Invalid migration password file' >&2; exit 2;; esac
case "$v05_collector_password" in ''|*[!0-9a-fA-F]*) printf '%s\n' 'Invalid collector password file' >&2; exit 2;; esac
case "$v05_publisher_password" in ''|*[!0-9a-fA-F]*) printf '%s\n' 'Invalid publisher password file' >&2; exit 2;; esac
if [ ! -f secrets/trade-log-publisher-database-url ]; then
  (umask 077; printf 'postgresql://trade_log_publisher:%s@postgres:5432/robotech\n' "$v05_publisher_password" > secrets/trade-log-publisher-database-url)
fi
if [ ! -f secrets/trade-log-collector-database-url ]; then
  (umask 077; printf 'postgresql://trade_log_collector:%s@postgres:5432/robotech\n' "$v05_collector_password" > secrets/trade-log-collector-database-url)
fi
if [ ! -f secrets/trade-log-database-url ]; then
  (umask 077; printf 'postgresql://trade_log_app:%s@postgres:5432/robotech\n' "$v05_app_password" > secrets/trade-log-database-url)
fi
if [ ! -f secrets/trade-log-migration-url ]; then
  (umask 077; printf 'postgresql://trade_log_migrator:%s@postgres:5432/robotech\n' "$v05_migration_password" > secrets/trade-log-migration-url)
fi
chmod 600 secrets/trade-log-publisher-password
chmod 644 secrets/trade-log-publisher-database-url secrets/webhook-token
chmod 600 secrets/trade-log-app-password secrets/trade-log-migrator-password secrets/trade-log-collector-password
chmod 644 secrets/trade-log-token secrets/postgres-password secrets/trade-log-database-url secrets/trade-log-migration-url secrets/trade-log-collector-database-url
# Stop writers before applying the new schema; preserve all containers' volumes.
docker compose stop query-api trade-log-query trade-collector trade-parser-publisher
docker compose up -d --wait --wait-timeout 90 postgres
# SQL is piped through stdin, passwords are not printed or put in command arguments.
docker compose exec -T postgres psql -q -v ON_ERROR_STOP=1 -U robotech_admin -d robotech <<SQL
SELECT format('CREATE ROLE trade_log_migrator LOGIN PASSWORD %L', '$v05_migration_password') WHERE NOT EXISTS (SELECT FROM pg_roles WHERE rolname='trade_log_migrator')\gexec
SELECT format('CREATE ROLE trade_log_app LOGIN PASSWORD %L', '$v05_app_password') WHERE NOT EXISTS (SELECT FROM pg_roles WHERE rolname='trade_log_app')\gexec
SELECT format('CREATE ROLE trade_log_collector LOGIN PASSWORD %L', '$v05_collector_password') WHERE NOT EXISTS (SELECT FROM pg_roles WHERE rolname='trade_log_collector')\gexec
SELECT format('CREATE ROLE trade_log_publisher LOGIN PASSWORD %L', '$v05_publisher_password') WHERE NOT EXISTS (SELECT FROM pg_roles WHERE rolname='trade_log_publisher')\gexec
GRANT CONNECT ON DATABASE robotech TO trade_log_publisher;
GRANT CONNECT ON DATABASE robotech TO trade_log_collector;
GRANT CONNECT, CREATE ON DATABASE robotech TO trade_log_migrator;
GRANT CREATE, USAGE ON SCHEMA public TO trade_log_migrator;
GRANT CONNECT ON DATABASE robotech TO trade_log_app;
SQL
unset v05_app_password v05_migration_password v05_collector_password v05_publisher_password
docker compose --profile init run --build --rm evidence-init
docker compose --profile init run --rm trade-log-migrate
