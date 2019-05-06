package entrypoint

import (
	"github.com/gorilla/mux"
)

func registerRoutes(r *mux.Router, serverIP, serverPort, publicKeyPath, privateKeyPath string, dev bool) {

	registerSuccessRoutes(r, serverIP, serverPort, publicKeyPath, privateKeyPath, dev)
	registerAuthFailureRoutes(r, serverIP, serverPort, publicKeyPath, privateKeyPath, dev)
	registerCertInvalidRoutes(r, serverIP, serverPort, publicKeyPath, privateKeyPath, dev)
	registerTokenInvalidRoutes(r, serverIP, serverPort, publicKeyPath, privateKeyPath, dev)
	registerCertMissingRoutes(r, serverIP, serverPort, publicKeyPath, privateKeyPath, dev)
	registerTokenMissingRoutes(r, serverIP, serverPort, publicKeyPath, privateKeyPath, dev)
	registerHealthzRoutes(r, serverIP, serverPort, publicKeyPath, privateKeyPath, dev)
}
