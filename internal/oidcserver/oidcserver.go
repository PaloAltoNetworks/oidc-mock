package oidcserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	jose "gopkg.in/square/go-jose.v2"
)

// location of the files used for signing and verification
const (
	privKeyPath = "app.rsa"     // openssl genrsa -out app.rsa
	pubKeyPath  = "app.rsa.pub" // openssl rsa -in app.rsa -pubout > app.rsa.pub
)

// NewOIDCServer returns ODIC handler
func NewOIDCServer() OIDCServer {

	return &oidcServer{
		rsa:   newRSAProcessor(pubKeyPath, privKeyPath),
		keyID: uuid.Must(uuid.NewV4()).String(),
	}
}

func (o *oidcServer) ProviderEndpoints(w http.ResponseWriter, r *http.Request) {

	fmt.Println("PROVIDER ENDPOINT")

	p := providerEndpoints{
		Issuer:      "https://192.168.100.1:6999",
		AuthURL:     "https://192.168.100.1:6999/auth",
		TokenURL:    "https://192.168.100.1:6999/token",
		UserInfoURL: "https://192.168.100.1:6999/userInfo",
		JWKSURL:     "https://192.168.100.1:6999/cert",
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (o *oidcServer) Authenticate(w http.ResponseWriter, r *http.Request) {

	fmt.Println("AUTH")

	state := r.URL.Query().Get("state")
	reqURI := r.URL.Query().Get("redirect_uri")

	req, err := http.NewRequest("GET", reqURI, nil)
	if err != nil {
		panic(err)
	}

	q := req.URL.Query()
	q.Add("state", state)
	q.Add("redirect_uri", reqURI)
	req.URL.RawQuery = q.Encode()

	http.Redirect(w, r, req.URL.String(), http.StatusTemporaryRedirect)
}

func (o *oidcServer) IssueToken(w http.ResponseWriter, r *http.Request) {

	fmt.Println("ISSUE")

	tokenExpiry := time.Now().AddDate(100, 0, 0).Unix()

	claims := jwt.MapClaims{
		"sub":  "1234567890",
		"iss":  "https://192.168.100.1:6999",
		"name": "Sibi",
		"exp":  tokenExpiry,
		"aud":  "450167263420-mrcsa5oue3jdm34a81e06gm3sk15t838.apps.googleusercontent.com",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = o.keyID

	idToken, err := token.SignedString(o.rsa.signKey())
	if err != nil {
		fmt.Println("SIGNING ERROR", err)
	}

	p := tokens{
		IDToken:     idToken,
		AccessToken: "https://192.168.100.1:6999/auth",
		TokenType:   "Bearer",
		ExpiresIn:   90000,
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (o *oidcServer) IssueCertificate(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Cert")

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
}

func (o *oidcServer) UserInfo(w http.ResponseWriter, r *http.Request) {
	// do nothing
}
