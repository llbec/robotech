param (
  [string]$ServiceName
)

if (-not $ServiceName) {
  Write-Host "Usage: .\gen_service_cert.ps1 <service-name>"
  exit 1
}

openssl genrsa -out "$ServiceName.key" 2048

@"
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
CN=$ServiceName

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = $ServiceName
DNS.2 = localhost
IP.1  = 127.0.0.1
"@ | Out-File "$ServiceName.cnf" -Encoding ascii

openssl req -new `
  -key "$ServiceName.key" `
  -out "$ServiceName.csr" `
  -config "$ServiceName.cnf"

openssl x509 -req `
  -in "$ServiceName.csr" `
  -CA ca.crt `
  -CAkey ca.key `
  -CAcreateserial `
  -out "$ServiceName.crt" `
  -days 365 `
  -sha256 `
  -extensions req_ext `
  -extfile "$ServiceName.cnf"

Remove-Item "$ServiceName.csr", "$ServiceName.cnf"

Write-Host "Certificate generated for $ServiceName"
