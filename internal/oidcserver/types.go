package oidcserver

import (
	"net/http"
	"time"
)

// OIDCServer exposes oidc server methods
type OIDCServer interface {
	ProviderEndpoints(w http.ResponseWriter, r *http.Request)
	Authenticate(w http.ResponseWriter, r *http.Request)
	IssueToken(w http.ResponseWriter, r *http.Request)
	IssueCertificate(w http.ResponseWriter, r *http.Request)
	UserInfo(w http.ResponseWriter, r *http.Request)
}

type oidcServer struct {
	rsa   *rsaProcessor
	keyID string
}

type providerEndpoints struct {
	Issuer      string `json:"issuer"`
	AuthURL     string `json:"authorization_endpoint"`
	TokenURL    string `json:"token_endpoint"`
	JWKSURL     string `json:"jwks_uri"`
	UserInfoURL string `json:"userinfo_endpoint"`
}

type tokens struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	IDToken      string        `json:"id_token"`
	TokenType    string        `json:"token_type"`
	ExpiresIn    time.Duration `json:"expires_in"`
}

type signingKeys struct {
	Keys []keys `json:"keys"`
}

type keys struct {
	Use string `json:"use"`
	Kid string `json:"kid"`
	E   string `json:"e"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	N   string `json:"n"`
}
