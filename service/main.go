package main

import (
	"ivory/persistence"
	"ivory/router"
	"ivory/service"
)

func main() {
	persistence.BoltDB.Build("cluster.db")
	persistence.File.Build()
	service.Secret.Create(persistence.BoltDB.Credential.GetEncryptedRef())
	router.Start()
}
