#!/usr/bin/env bash

set -euo pipefail

# 导入参数：按需修改这三个值。
address="0x45d26f28196d226497130c4bac709d808fed4029"
from_time="2026-09-01T00:00:00Z"
to_time="2026-09-02T00:00:00Z"

script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
repo_dir="$(cd -- "${script_dir}/.." && pwd)"
env_file="${ENV_FILE:-${repo_dir}/.env}"

if [[ ! -f "${env_file}" ]]; then
  echo "配置文件不存在：${env_file}" >&2
  echo "请创建该文件并设置 DATABASE_URL、NANSEN_API_KEY，可选设置 RUST_LOG。" >&2
  exit 2
fi

set -a
# shellcheck disable=SC1090
source "${env_file}"
set +a

: "${DATABASE_URL:?DATABASE_URL 未在 ${env_file} 中设置}"
: "${NANSEN_API_KEY:?NANSEN_API_KEY 未在 ${env_file} 中设置}"
export RUST_LOG="${RUST_LOG:-info}"

cd "${repo_dir}"

exec cargo run -p trade-log-import -- \
  --address "${address}" \
  --from "${from_time}" \
  --to "${to_time}"
