package oidcserver

import "fmt"

func generateProviderURLs(ip, port string) providerEndpoints {

	return providerEndpoints{
		Issuer:      generateCompleteURL(ip, port, ""),
		AuthURL:     generateCompleteURL(ip, port, "/auth"),
		TokenURL:    generateCompleteURL(ip, port, "/token"),
		UserInfoURL: generateCompleteURL(ip, port, "/userInfo"),
		JWKSURL:     generateCompleteURL(ip, port, "/cert"),
	}
}

func generateCompleteURL(ip, port, endpoint string) string {

	return fmt.Sprintf("https://%s%s%s", ip, port, endpoint)
}
