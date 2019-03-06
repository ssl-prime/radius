package main

import (
	radius "radius/api/service"

	"github.com/rightjoin/aqua"
)

// main function
func main() {
	server := aqua.NewRestServer()
	server.AddService(&radius.RepoInfo{})

	server.Run()
}
