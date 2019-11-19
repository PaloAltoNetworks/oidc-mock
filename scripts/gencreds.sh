#!/bin/bash -e

usage()
{
    echo "usage: gencreds.sh [--dns hello.com]... [--ip 127.0.0.1]... --force"
}

# CA_SAN can be set externally and imported as well
# check if dns or ip is supplied on the command line
while [ "$1" != "" ]; do
  case $1 in
    --dns | --ip )
        if [ $# -gt 1 ]; then
            CA_SAN="${CA_SAN} $1 $2"
        else
            usage
            exit
        fi
        ;;
    -f | --force ) FORCE=1 ;;
    -h | --help )
        usage
        exit
        ;;
  esac
  shift
done

if [ -z "${CA_SAN}" ]; then
  usage
  exit
fi

# create creds folder
if [ "${FORCE}" = "1" ]; then
  echo "Removing .data"
  rm -rf .data
fi
set +e

if ! mkdir .data ; then
  echo "Use --force at the end to remove .data permanently"
  exit 1
fi
set -e

# generate self signed ca and key
# shellcheck disable=SC2086
tg cert --name system --org oidc.com --common-name oidc-mock --pass oidc ${CA_SAN} --is-ca
mv system-cert.pem .data/system.crt
# decrypt the system key
tg decrypt --key system-key.pem --pass oidc > .data/system.key
rm -rf system-key.pem
# generate ca signed server cert and key
# shellcheck disable=SC2086
tg cert --name server --org oidc.com --common-name oidc-mock --pass oidc ${CA_SAN} --auth-server --signing-cert .data/system.crt --signing-cert-key .data/system.key --signing-cert-key-pass oidc
mv server-cert.pem .data/server.crt
# decrypt the server key
tg decrypt --key server-key.pem --pass oidc > .data/server.key
rm -rf server-key.pem

# generate public,private key for token verification
openssl genrsa -out oidc.rsa
openssl rsa -in oidc.rsa -pubout > oidc.rsa.pub
mv oidc.rsa .data/oidc.rsa
mv oidc.rsa.pub .data/oidc.rsa.pub
