package main

import (
	"go-jwt-mysql/config"
	"go-jwt-mysql/route"
)

var err error

func main() {
	config.Init()

	r := route.SetupRouter()
	//running
	r.Run()

	defer config.CloseDb()
}
