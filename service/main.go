package main

import (
	"ivory/persistence"
	"ivory/router"
)

func main() {
	persistence.Database.Build("cluster.db")
	persistence.File.Build()
	router.Start()
}
