package main

import (
	"ivory/src/app"
)

func main() {
	context := app.NewContext()
	app.NewRouter(context)
}
