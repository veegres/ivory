package main

import (
	"ivory/persistence"
	"ivory/router"
	"ivory/service"
)

func main() {
	persistence.BoltDB.Build("cluster.db")
	persistence.File.Build()
	ref, err := persistence.BoltDB.Credential.GetEncryptedRef()
	if err != nil {
		panic(err)
	}
	service.Secret.Create(ref)
	router.Start()
}
