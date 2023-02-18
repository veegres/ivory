package main

import (
	"ivory/app"
)

func main() {
	di := app.NewDI()
	app.NewRouter(di)
}
