#!/bin/bash

# 1. Generate server priv/pub key pair
openssl genpkey -out priv.key -outpubkey pub.key -algorithm RSA -pkeyopt rsa_keygen_bits:4096

# 2. Generate root CA priv key
openssl genpkey -out ca-priv.key -outpubkey ca-pub.key -algorithm RSA -pkeyopt rsa_keygen_bits:4096

# 3. Generate root CA certificate
openssl x509 -new -key ca-priv.key -subj /DC=CA/DC=RootCA/DC=RootCA/UID=111111+CN=RootCA -out ca.cert

# 4. Generate server CSR
openssl req -new -key priv.key -subj /DC=ServCert/DC=ServCert/DC=ServCert/UID=222222+CN=ServCert -out serv-cert.req

# 5. Sign server CSR with root CA
openssl req -in serv-cert.req -x509 -copy_extensions copyall -CA ca.cert -CAkey ca-priv.key -out serv.cert


