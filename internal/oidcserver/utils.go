package oidcserver

import (
	"fmt"
	"path"
)

func generateProviderURLs(serverFlow ServerFlowType, ip, port string, devMode bool) providerEndpoints {

	return providerEndpoints{
		Issuer:      generateCompleteURL(serverFlow, ip, port, "", devMode),
		AuthURL:     generateCompleteURL(serverFlow, ip, port, "/auth", devMode),
		TokenURL:    generateCompleteURL(serverFlow, ip, port, "/token", devMode),
		UserInfoURL: generateCompleteURL(serverFlow, ip, port, "/userInfo", devMode),
		JWKSURL:     generateCompleteURL(serverFlow, ip, port, "/cert", devMode),
	}
}

func generateCompleteURL(serverFlow ServerFlowType, ip, port, endpoint string, devMode bool) string {

	if endpoint == "" {
		endpoint = "/"
	}

	switch serverFlow {
	case ServerFlowTypeAuthFailure:
		endpoint = path.Join(endpoint, AuthFailure)
	case ServerFlowTypeInvalidCert:
		endpoint = path.Join(endpoint, CertInvalid)
	case ServerFlowTypeInvalidToken:
		endpoint = path.Join(endpoint, TokenInvalid)
	}

	if devMode {
		return fmt.Sprintf("https://%s%s%s", ip, port, endpoint)
	}

	return fmt.Sprintf("https://%s%s", ip, endpoint)
}
