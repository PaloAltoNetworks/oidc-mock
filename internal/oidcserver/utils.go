package oidcserver

import (
	"fmt"
	"path"
)

func generateProviderURLs(serverFlow ServerFlowType, ip, port string) providerEndpoints {

	return providerEndpoints{
		Issuer:      generateCompleteURL(serverFlow, ip, port, ""),
		AuthURL:     generateCompleteURL(serverFlow, ip, port, "/auth"),
		TokenURL:    generateCompleteURL(serverFlow, ip, port, "/token"),
		UserInfoURL: generateCompleteURL(serverFlow, ip, port, "/userInfo"),
		JWKSURL:     generateCompleteURL(serverFlow, ip, port, "/cert"),
	}
}

func generateCompleteURL(serverFlow ServerFlowType, ip, port, endpoint string) string {

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

	return fmt.Sprintf("https://%s%s%s", ip, port, endpoint)
}
