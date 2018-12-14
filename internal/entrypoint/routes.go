package entrypoint

import (
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"go.aporeto.io/oidc-mock/internal/oidcserver"
)

func registerSuccessRoutes(r *mux.Router, serverIP, serverPort, publicKeyPath, privateKeyPath string, dev bool) {

	oidc := oidcserver.NewOIDCServer(oidcserver.ServerFlowTypeSuccess, serverIP, serverPort, publicKeyPath, privateKeyPath, dev)

	r.HandleFunc("/.well-known/openid-configuration", oidc.ProviderEndpoints).Methods(http.MethodGet)
	r.HandleFunc("/auth", oidc.Authenticate).Methods(http.MethodGet)
	r.HandleFunc("/userInfo", oidc.UserInfo).Methods(http.MethodGet)
	r.HandleFunc("/token", oidc.IssueToken).Methods(http.MethodPost)
	r.HandleFunc("/cert", oidc.IssueCertificate).Methods(http.MethodGet)
}

func registerAuthFailureRoutes(r *mux.Router, serverIP, serverPort, publicKeyPath, privateKeyPath string, dev bool) {

	oidc := oidcserver.NewOIDCServer(oidcserver.ServerFlowTypeAuthFailure, serverIP, serverPort, publicKeyPath, privateKeyPath, dev)

	r.HandleFunc(path.Join("/"+oidcserver.AuthFailure, ".well-known/openid-configuration"), oidc.ProviderEndpoints).Methods(http.MethodGet)
	r.HandleFunc(path.Join("/auth", oidcserver.AuthFailure), oidc.Authenticate).Methods(http.MethodGet)
	r.HandleFunc(path.Join("/userInfo", oidcserver.AuthFailure), oidc.UserInfo).Methods(http.MethodGet)
	r.HandleFunc(path.Join("/token", oidcserver.AuthFailure), oidc.IssueToken).Methods(http.MethodPost)
	r.HandleFunc(path.Join("/cert", oidcserver.AuthFailure), oidc.IssueCertificate).Methods(http.MethodGet)
}

func registerTokenInvalidRoutes(r *mux.Router, serverIP, serverPort, publicKeyPath, privateKeyPath string, dev bool) {

	oidc := oidcserver.NewOIDCServer(oidcserver.ServerFlowTypeInvalidToken, serverIP, serverPort, publicKeyPath, privateKeyPath, dev)

	r.HandleFunc(path.Join("/"+oidcserver.TokenInvalid, ".well-known/openid-configuration"), oidc.ProviderEndpoints).Methods(http.MethodGet)
	r.HandleFunc(path.Join("/auth", oidcserver.TokenInvalid), oidc.Authenticate).Methods(http.MethodGet)
	r.HandleFunc(path.Join("/userInfo", oidcserver.TokenInvalid), oidc.UserInfo).Methods(http.MethodGet)
	r.HandleFunc(path.Join("/token", oidcserver.TokenInvalid), oidc.IssueToken).Methods(http.MethodPost)
	r.HandleFunc(path.Join("/cert", oidcserver.TokenInvalid), oidc.IssueCertificate).Methods(http.MethodGet)
}

func registerCertInvalidRoutes(r *mux.Router, serverIP, serverPort, publicKeyPath, privateKeyPath string, dev bool) {

	oidc := oidcserver.NewOIDCServer(oidcserver.ServerFlowTypeInvalidCert, serverIP, serverPort, publicKeyPath, privateKeyPath, dev)

	r.HandleFunc(path.Join("/"+oidcserver.CertInvalid, ".well-known/openid-configuration"), oidc.ProviderEndpoints).Methods(http.MethodGet)
	r.HandleFunc(path.Join("/auth", oidcserver.CertInvalid), oidc.Authenticate).Methods(http.MethodGet)
	r.HandleFunc(path.Join("/userInfo", oidcserver.CertInvalid), oidc.UserInfo).Methods(http.MethodGet)
	r.HandleFunc(path.Join("/token", oidcserver.CertInvalid), oidc.IssueToken).Methods(http.MethodPost)
	r.HandleFunc(path.Join("/cert", oidcserver.CertInvalid), oidc.IssueCertificate).Methods(http.MethodGet)
}

func registerCertMissingRoutes(r *mux.Router, serverIP, serverPort, publicKeyPath, privateKeyPath string, dev bool) {

	oidc := oidcserver.NewOIDCServer(oidcserver.ServerFlowTypeMissingCert, serverIP, serverPort, publicKeyPath, privateKeyPath, dev)

	r.HandleFunc(path.Join("/"+oidcserver.CertMissing, ".well-known/openid-configuration"), oidc.ProviderEndpoints).Methods(http.MethodGet)
	r.HandleFunc(path.Join("/auth", oidcserver.CertMissing), oidc.Authenticate).Methods(http.MethodGet)
	r.HandleFunc(path.Join("/userInfo", oidcserver.CertMissing), oidc.UserInfo).Methods(http.MethodGet)
	r.HandleFunc(path.Join("/token", oidcserver.CertMissing), oidc.IssueToken).Methods(http.MethodPost)
	r.HandleFunc(path.Join("/cert", oidcserver.CertMissing), oidc.IssueCertificate).Methods(http.MethodGet)
}

func registerTokenMissingRoutes(r *mux.Router, serverIP, serverPort, publicKeyPath, privateKeyPath string, dev bool) {

	oidc := oidcserver.NewOIDCServer(oidcserver.ServerFlowTypeMissingToken, serverIP, serverPort, publicKeyPath, privateKeyPath, dev)

	r.HandleFunc(path.Join("/"+oidcserver.TokenMissing, ".well-known/openid-configuration"), oidc.ProviderEndpoints).Methods(http.MethodGet)
	r.HandleFunc(path.Join("/auth", oidcserver.TokenMissing), oidc.Authenticate).Methods(http.MethodGet)
	r.HandleFunc(path.Join("/userInfo", oidcserver.TokenMissing), oidc.UserInfo).Methods(http.MethodGet)
	r.HandleFunc(path.Join("/token", oidcserver.TokenMissing), oidc.IssueToken).Methods(http.MethodPost)
	r.HandleFunc(path.Join("/cert", oidcserver.TokenMissing), oidc.IssueCertificate).Methods(http.MethodGet)
}
