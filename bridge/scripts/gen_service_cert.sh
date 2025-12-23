#!/usr/bin/env bash
set -e

SERVICE_NAME=$1
DAYS=365

if [ -z "$SERVICE_NAME" ]; then
  echo "Usage: $0 <service-name>"
  exit 1
fi

echo "==> Generating private key for ${SERVICE_NAME}"
openssl genrsa -out ${SERVICE_NAME}.key 2048

cat > ${SERVICE_NAME}.cnf <<EOF
[ req ]
default_bits       = 2048
prompt             = no
default_md         = sha256
distinguished_name = dn
req_extensions     = req_ext

[ dn ]
C=CN
O=BridgeSystem
OU=Service
CN=${SERVICE_NAME}

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = ${SERVICE_NAME}
DNS.2 = ${SERVICE_NAME}.default
DNS.3 = localhost
IP.1  = 127.0.0.1
EOF

echo "==> Generating CSR"
openssl req -new \
  -key ${SERVICE_NAME}.key \
  -out ${SERVICE_NAME}.csr \
  -config ${SERVICE_NAME}.cnf

echo "==> Signing certificate with CA"
openssl x509 -req \
  -in ${SERVICE_NAME}.csr \
  -CA ca.crt \
  -CAkey ca.key \
  -CAcreateserial \
  -out ${SERVICE_NAME}.crt \
  -days ${DAYS} \
  -sha256 \
  -extensions req_ext \
  -extfile ${SERVICE_NAME}.cnf

rm -f ${SERVICE_NAME}.csr ${SERVICE_NAME}.cnf

echo "==> Generated:"
ls -l ${SERVICE_NAME}.key ${SERVICE_NAME}.crt


#./gen_service_cert.sh bridge
#./gen_service_cert.sh txstore
#./gen_service_cert.sh report-daily
