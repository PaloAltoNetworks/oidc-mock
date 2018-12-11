#!/bin/bash -e

CA_SAN=${CA_SAN:-$1}

[ ! -z ${CA_SAN} ] || (echo "Missing ip" && exit 1)

# generate signing cert and key
tg cert --name system --org oidc.com --common-name oidc-mock --ip ${CA_SAN} --auth-server
mv system-cert.pem .data/system.crt
mv system-key.pem .data/system.key

# generate public,private key
openssl genrsa -out oidc.rsa
openssl rsa -in oidc.rsa -pubout > oidc.rsa.pub
mv oidc.rsa .data/oidc.rsa
mv oidc.rsa.pub .data/oidc.rsa.pub
