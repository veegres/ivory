package main

import (
	"ivory/database"
	"ivory/router"
)

func main() {
	database.Start()
	router.Start()
}
