package main

import (
	"ivory/persistence"
	"ivory/router"
	"ivory/service"
)

func main() {
	persistence.Database.Build("cluster.db")
	persistence.File.Build()
	service.Secret.Create(persistence.Database.Credential.GetEncryptedRef())
	router.Start()
}
