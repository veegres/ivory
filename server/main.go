package main

import (
	core2 "ivory/core"
)

func main() {
	context := core2.NewContext()
	core2.NewRouter(context)
}
