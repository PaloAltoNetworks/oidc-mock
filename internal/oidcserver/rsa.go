package oidcserver

import (
	"crypto/rsa"
	"io/ioutil"

	jwt "github.com/dgrijalva/jwt-go"
)

type rsaProcessor struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func newRSAProcessor(publicKeyPath, privateKeyPath string) *rsaProcessor {

	publicKey, privateKey := retrieveKeys(publicKeyPath, privateKeyPath)

	return &rsaProcessor{
		publicKey:  publicKey,
		privateKey: privateKey,
	}
}

func (r *rsaProcessor) signKey() *rsa.PrivateKey {

	return r.privateKey
}

func (r *rsaProcessor) verifyKey() *rsa.PublicKey {

	return r.publicKey
}

func retrieveKeys(privateKeyPath, publicKeyPath string) (*rsa.PublicKey, *rsa.PrivateKey) {

	signBytes, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		panic(err)
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		panic(err)
	}

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		panic(err)
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		panic(err)
	}

	return verifyKey, signKey
}
