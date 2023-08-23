#!/bin/sh

# Generate a private key
openssl genrsa -out server.key 2048

# Generate a certificate signing request (CSR)
openssl req -new -key server.key -out server.csr

# Generate a self-signed certificate using the CSR and private key
openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt

# Clean up the CSR
rm server.csr