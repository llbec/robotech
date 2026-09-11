#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
repo_dir="$(cd -- "${script_dir}/.." && pwd)"
env_file="${ENV_FILE:-${repo_dir}/.env}"
import_config_file="${IMPORT_CONFIG_FILE:-${repo_dir}/.run-v0.000.env}"

if [[ ! -f "${env_file}" ]]; then
  echo "配置文件不存在：${env_file}" >&2
  echo "请创建该文件并设置 DATABASE_URL、NANSEN_API_KEY，可选设置 RUST_LOG。" >&2
  exit 2
fi

if [[ ! -f "${import_config_file}" ]]; then
  echo "导入配置文件不存在：${import_config_file}" >&2
  echo "请设置 IMPORT_ADDRESS、IMPORT_FROM、IMPORT_TO。" >&2
  exit 2
fi

set -a
# shellcheck disable=SC1090
source "${env_file}"
# shellcheck disable=SC1090
source "${import_config_file}"
set +a

: "${DATABASE_URL:?DATABASE_URL 未在 ${env_file} 中设置}"
: "${NANSEN_API_KEY:?NANSEN_API_KEY 未在 ${env_file} 中设置}"
: "${IMPORT_ADDRESS:?IMPORT_ADDRESS 未在 ${import_config_file} 中设置}"
: "${IMPORT_FROM:?IMPORT_FROM 未在 ${import_config_file} 中设置}"
: "${IMPORT_TO:?IMPORT_TO 未在 ${import_config_file} 中设置}"
export RUST_LOG="${RUST_LOG:-info}"

cd "${repo_dir}"

exec cargo run -p trade-log-import -- \
  --address "${IMPORT_ADDRESS}" \
  --from "${IMPORT_FROM}" \
  --to "${IMPORT_TO}"
