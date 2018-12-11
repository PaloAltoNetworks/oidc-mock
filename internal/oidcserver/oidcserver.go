package oidcserver

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	jose "gopkg.in/square/go-jose.v2"
)

// NewOIDCServer returns ODIC handler
func NewOIDCServer(serverIP, serverPort, publicKeyPath, privateKeyPath string) OIDCServer {

	return &oidcServer{
		rsa:        newRSAProcessor(publicKeyPath, privateKeyPath),
		keyID:      uuid.Must(uuid.NewV4()).String(),
		serverIP:   serverIP,
		serverPort: serverPort,
	}
}

// ProviderEndpoints returns provider urls for lib
func (o *oidcServer) ProviderEndpoints(w http.ResponseWriter, r *http.Request) {

	zap.L().Debug("ProviderEndpoints called")

	providerURLs := generateProviderURLs(o.serverIP, o.serverPort)

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(providerURLs)
}

// Authenticate is a mock call which by default redirects to redirect_uri given in request
// NOTE: There is NO authentication is done here
func (o *oidcServer) Authenticate(w http.ResponseWriter, r *http.Request) {

	zap.L().Debug("Authenticate called")

	state := r.URL.Query().Get("state")
	redURI := r.URL.Query().Get("redirect_uri")

	reqURI, err := url.ParseRequestURI(redURI)
	if err != nil {
		zap.L().Error("Unable to parse redirect uri", zap.Error(err))
		return
	}

	q := reqURI.Query()
	q.Add("state", state)
	q.Add("redirect_uri", redURI)
	reqURI.RawQuery = q.Encode()

	http.Redirect(w, r, reqURI.String(), http.StatusTemporaryRedirect)
}

// IssueToken issues JWT token
func (o *oidcServer) IssueToken(w http.ResponseWriter, r *http.Request) {

	zap.L().Debug("IssueToken called")

	tokenExpiry := time.Now().AddDate(100, 0, 0).Unix()

	claims := jwt.MapClaims{
		"sub":  "1234567890",
		"iss":  generateCompleteURL(o.serverIP, o.serverPort, ""),
		"name": "oidc-mock",
		"exp":  tokenExpiry,
		"aud":  "apps.oidcmock.com",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	// NOTE: This is important as the library matches this keyID with the public key
	token.Header["kid"] = o.keyID

	idToken, err := token.SignedString(o.rsa.signKey())
	if err != nil {
		zap.L().Error("Unable to sign JWT", zap.Error(err))
		return
	}

	p := tokens{
		IDToken:     idToken,
		AccessToken: "notoken",
		TokenType:   "Bearer",
		ExpiresIn:   90000,
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)

	zap.L().Debug("Token issued")
}

// IssueToken issues public certificate used to sign JWT
func (o *oidcServer) IssueCertificate(w http.ResponseWriter, r *http.Request) {

	zap.L().Debug("IssueCertificate called")

	jwk := jose.JSONWebKey{
		Key:       o.rsa.verifyKey(),
		KeyID:     o.keyID,
		Use:       "sig",
		Algorithm: "RS256",
	}

	jwks := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{jwk},
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwks)

	zap.L().Debug("Certificate issued")
}

// UserInfo ...
func (o *oidcServer) UserInfo(w http.ResponseWriter, r *http.Request) {
	zap.L().Debug("Userinfo called")
	// Do nothing
}