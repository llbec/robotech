#!/bin/sh
set -eu
if ! command -v python3 >/dev/null 2>&1; then
  echo '[FAIL] 缺少 python3（仅使用标准库）' >&2
  exit 2
fi
script_directory=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
exec python3 "$script_directory/verify-version.py" --version v0.2 "$@"
