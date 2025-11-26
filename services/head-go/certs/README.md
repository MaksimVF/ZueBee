

# Certificates for mTLS

This directory contains certificates for mTLS (mutual TLS) authentication.

## Generating Self-Signed Certificates

To generate self-signed certificates for development, you can use OpenSSL:

```bash
# Generate CA key and certificate
openssl genpkey -algorithm RSA -out ca.key
openssl req -x509 -new -nodes -key ca.key -sha256 -days 365 -out ca.crt -subj "/CN=My CA"

# Generate server key and certificate signing request (CSR)
openssl genpkey -algorithm RSA -out head.key
openssl req -new -key head.key -out head.csr -subj "/CN=head-service"

# Sign the server certificate with the CA
openssl x509 -req -in head.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out head.crt -days 365 -sha256
```

## Files

- `ca.crt`: CA certificate
- `head.crt`: Server certificate
- `head.key`: Server private key

## Usage

These certificates are used for mTLS authentication in the head-service. Update your configuration to point to these files.

