#!/bin/bash -e

CA_SAN=${CA_SAN:-$2}
FORCE=${FORCE:-$3}

usage()
{
    echo "usage: gencreds.sh [--dns] ['hello.com'] --force"
}

DNS=0
IP=0
# check if its dns or ip
while [ "$1" != "" ]; do
    case $1 in
            --dns )       DNS=1
                          ;;
             --ip )       IP=1
                          ;;
      -h | --help )       usage
                          exit
  esac
  shift
done

# validate and populate field
if [ "${DNS}" = "1" ]; then
  [ ! -z ${CA_SAN} ] || (echo "Missing dns" && exit 1)
    echo "Using DNS ${CA_SAN}"
    CA_SAN="--dns ${CA_SAN}"
elif [ "${IP}" = "1" ]; then
  [ ! -z ${CA_SAN} ] || (echo "Missing ip" && exit 1)
  echo "Using IP ${CA_SAN}"
  CA_SAN="--ip ${CA_SAN}"
fi

# create creds folder
if [ "${FORCE}" = "--force" ]; then
  echo "Removing .data"
  rm -rf .data
fi
set +e
mkdir .data
if [ $? != 0 ]; then
  echo "Use --force at the end to remove .data permanently"
  exit 1
fi
set -e

# generate self signed ca and key
tg cert --name system --org oidc.com --common-name oidc-mock --pass oidc  ${CA_SAN} --is-ca
mv system-cert.pem .data/system.crt
# decrypt the system key
tg decrypt --key system-key.pem --pass oidc > .data/system.key
rm -rf system-key.pem
# generate ca signed server cert and key
tg cert --name server --org oidc.com --common-name oidc-mock --pass oidc  ${CA_SAN} --auth-server --signing-cert .data/system.crt --signing-cert-key .data/system.key --signing-cert-key-pass oidc
mv server-cert.pem .data/server.crt
# decrypt the server key
tg decrypt --key server-key.pem --pass oidc > .data/server.key
rm -rf server-key.pem

# generate public,private key for token verification
openssl genrsa -out oidc.rsa
openssl rsa -in oidc.rsa -pubout > oidc.rsa.pub
mv oidc.rsa .data/oidc.rsa
mv oidc.rsa.pub .data/oidc.rsa.pub
