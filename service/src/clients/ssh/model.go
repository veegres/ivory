package ssh

import "crypto/ed25519"

type Connection struct {
	Host       string
	Port       int
	Username   string
	PrivateKey ed25519.PrivateKey
}

type CommandResult struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exitCode"`
}
