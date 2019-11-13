package oidcserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	jose "gopkg.in/square/go-jose.v2"
)

// NewOIDCServer returns ODIC handler
func NewOIDCServer(serverFlow ServerFlowType, serverIP, serverPort, publicKeyPath, privateKeyPath string, devMode bool) OIDCServer {

	return &oidcServer{
		rsa:        newRSAProcessor(serverFlow, publicKeyPath, privateKeyPath),
		keyID:      uuid.Must(uuid.NewV4()).String(),
		serverIP:   serverIP,
		serverPort: serverPort,
		serverFlow: serverFlow,
		devMode:    devMode,
	}
}

// ProviderEndpoints returns provider urls for lib
func (o *oidcServer) ProviderEndpoints(w http.ResponseWriter, r *http.Request) {

	o.Lock()
	defer o.Unlock()

	zap.L().Debug("Discovering Endpoints")

	providerURLs := o.generateProviderURLs()

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(providerURLs)
}

// Authenticate is a mock call which by default redirects to redirect_uri given in request
// NOTE: There is NO authentication is done here
func (o *oidcServer) Authenticate(w http.ResponseWriter, r *http.Request) {

	o.Lock()
	defer o.Unlock()

	zap.L().Debug("Authenticating")

	if o.serverFlow == ServerFlowTypeAuthFailure {
		http.Error(w, "Authentication failure", http.StatusUnauthorized)
		zap.L().Warn("Authentication failure", zap.Reflect("type", o.serverFlow))
		return
	}

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

	o.Lock()
	defer o.Unlock()

	zap.L().Debug("Issuing Token")

	tokenExpiry := time.Now().AddDate(100, 0, 0).Unix()

	claims := jwt.MapClaims{
		"sub":            "1234567890",
		"iss":            o.generateCompleteURL(""),
		"name":           "oidc-mock",
		"exp":            tokenExpiry,
		"aud":            "abcd1234.apps.oidcmock.com",
		"email":          "oidc-mock@example.com",
		"email_verified": true,
		"groups":         []string{"test", "dev"},
		"enabled":        true,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	// NOTE: This is important as the library matches this keyID with the public key
	token.Header["kid"] = o.keyID

	idToken, err := token.SignedString(o.rsa.signKey())
	if err != nil {
		zap.L().Error("Unable to sign JWT", zap.Error(err))
		return
	}

	if o.serverFlow == ServerFlowTypeMissingToken {
		idToken = ""
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

	o.Lock()
	defer o.Unlock()

	zap.L().Debug("Issuing Certificate")

	var verifyKey interface{}
	switch o.serverFlow {
	case ServerFlowTypeInvalidCert:
		verifyKey = []byte("invalidKey")
	case ServerFlowTypeMissingCert:
		verifyKey = []byte("")
	default:
		verifyKey = o.rsa.verifyKey()
	}

	jwk := jose.JSONWebKey{
		Key:       verifyKey,
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

	o.Lock()
	defer o.Unlock()

	zap.L().Debug("Userinfo called")
	// Do nothing
}

func (o *oidcServer) generateProviderURLs() providerEndpoints {

	return providerEndpoints{
		Issuer:      o.generateCompleteURL(""),
		AuthURL:     o.generateCompleteURL("/auth"),
		TokenURL:    o.generateCompleteURL("/token"),
		UserInfoURL: o.generateCompleteURL("/userInfo"),
		JWKSURL:     o.generateCompleteURL("/cert"),
	}
}

func (o *oidcServer) generateCompleteURL(endpoint string) string {

	if endpoint == "" && o.serverFlow != ServerFlowTypeSuccess {
		endpoint = "/"
	}

	switch o.serverFlow {
	case ServerFlowTypeAuthFailure:
		endpoint = path.Join(endpoint, AuthFailure)
	case ServerFlowTypeInvalidCert:
		endpoint = path.Join(endpoint, CertInvalid)
	case ServerFlowTypeInvalidToken:
		endpoint = path.Join(endpoint, TokenInvalid)
	case ServerFlowTypeMissingCert:
		endpoint = path.Join(endpoint, CertMissing)
	case ServerFlowTypeMissingToken:
		endpoint = path.Join(endpoint, TokenMissing)
	}

	if o.devMode {
		return fmt.Sprintf("https://%s%s%s", o.serverIP, o.serverPort, endpoint)
	}

	return fmt.Sprintf("https://%s%s", o.serverIP, endpoint)
}
