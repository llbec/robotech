#!/bin/sh
set -eu
ROBOTECH_TEST_DATABASE_URL=$(cat /run/secrets/fixture-database-url)
export ROBOTECH_TEST_DATABASE_URL
exec /app/spot-fixtures server_spot_suite --ignored --nocapture --test-threads=1
