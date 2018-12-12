#!/bin/bash -e

CA_SAN=${CA_SAN:-$1}

[ ! -z ${CA_SAN} ] || (echo "Missing ip" && exit 1)

# create creds folder
mkdir .data

# generate self signed ca
tg cert --name system --org oidc.com --common-name oidc-mock --pass oidc  --ip ${CA_SAN} --is-ca
mv system-cert.pem .data/system.crt
# decrypt the system key
tg decrypt --key system-key.pem --pass oidc > .data/system.key
rm -rf system-key.pem
# generate ca signed server cert and key
tg cert --name server --org oidc.com --common-name oidc-mock --pass oidc  --ip ${CA_SAN} --auth-server --signing-cert .data/system.crt --signing-cert-key .data/system.key --signing-cert-key-pass oidc
mv server-cert.pem .data/server.crt
# decrypt the server key
tg decrypt --key server-key.pem --pass oidc > .data/server.key
rm -rf server-key.pem

# generate public,private key for token verification
openssl genrsa -out oidc.rsa
openssl rsa -in oidc.rsa -pubout > oidc.rsa.pub
mv oidc.rsa .data/oidc.rsa
mv oidc.rsa.pub .data/oidc.rsa.pub
