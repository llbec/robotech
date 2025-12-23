#!/usr/bin/env bash
set -e

CA_NAME="bridge-ca"
DAYS=3650

echo "==> Generating CA private key"
openssl genrsa -out ca.key 4096

echo "==> Generating CA certificate"
openssl req -x509 -new -nodes \
  -key ca.key \
  -sha256 \
  -days ${DAYS} \
  -out ca.crt \
  -subj "/C=CN/O=BridgeSystem/OU=Infra/CN=${CA_NAME}"

echo "==> CA generated:"
ls -l ca.key ca.crt
