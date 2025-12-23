$CAName = "bridge-ca"
$Days = 3650

openssl genrsa -out ca.key 4096

openssl req -x509 -new -nodes `
  -key ca.key `
  -sha256 `
  -days $Days `
  -out ca.crt `
  -subj "/C=CN/O=BridgeSystem/OU=Infra/CN=$CAName"

Write-Host "CA generated"
