#!/bin/sh
set -eu
# Incremental schema upgrade; preserve existing data and secrets.
sh scripts/init-v0.6.sh
printf '%s\n' 'v0.7 initialized; start with docker compose up --build -d. Spot publication remains opt-in.'
