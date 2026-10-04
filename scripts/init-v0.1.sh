#!/bin/sh
set -eu
# Run from the repository root. Preserve an existing deployment credential.
mkdir -p secrets
chmod 700 secrets
if [ ! -f secrets/trade-log-token ]; then
  (umask 077; openssl rand -hex 32 > secrets/trade-log-token)
fi
# Container UID 10001 must be able to read the single mounted file. The host
# directory stays owner-only; the token is excluded from Git and build contexts.
chmod 644 secrets/trade-log-token
docker compose --profile init run --build --rm evidence-init
