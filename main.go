package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.aporeto.io/oidc-mock/internal/oidcserver"
)

func main() {

	r := mux.NewRouter()
	oidc := oidcserver.NewOIDCServer()

	r.HandleFunc("/.well-known/openid-configuration", oidc.ProviderEndpoints).Methods(http.MethodGet)
	r.HandleFunc("/auth", oidc.Authenticate).Methods(http.MethodGet)
	r.HandleFunc("/userInfo", oidc.UserInfo).Methods(http.MethodGet)
	r.HandleFunc("/token", oidc.IssueToken).Methods(http.MethodPost)
	r.HandleFunc("/cert", oidc.IssueCertificate).Methods(http.MethodGet)

	fmt.Println(http.ListenAndServeTLS(":6999", "selfsigned.crt", "selfsigned.key", r))
}
