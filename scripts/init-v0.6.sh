#!/bin/sh
set -eu
# Preserve existing data and credentials. v0.5 initializer runs the current migrator.
mkdir -p secrets
if [ ! -f secrets/watch-management-token ]; then
  (umask 077; openssl rand -hex 32 > secrets/watch-management-token)
fi
chmod 644 secrets/watch-management-token
sh scripts/init-v0.5.sh
printf '%s\n' 'v0.6 migration complete. Start applications with docker compose up --build -d; original checkpoints are imported without reset.'
