package oidcserver

import (
	"crypto/rsa"
	"io/ioutil"

	jwt "github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
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

func retrieveKeys(publicKeyPath, privateKeyPath string) (*rsa.PublicKey, *rsa.PrivateKey) {

	verifyBytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		zap.L().Fatal("Unable to read public key path",
			zap.String("path", publicKeyPath),
			zap.Error(err))
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		zap.L().Fatal("Unable to parse public key", zap.Error(err))
	}

	signBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		zap.L().Fatal("Unable to read private key path",
			zap.String("path", privateKeyPath),
			zap.Error(err))
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		zap.L().Fatal("Unable to parse private key", zap.Error(err))
	}

	return verifyKey, signKey
}
