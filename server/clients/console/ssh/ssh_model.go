package ssh

import "crypto/ed25519"

type Connection struct {
	Host       string
	Port       int
	Username   string
	Password   string
	PrivateKey *ed25519.PrivateKey
}
