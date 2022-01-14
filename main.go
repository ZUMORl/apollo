package main

import (
	"flag"

	"github.com/apollo/server"
)

func main() {
	var origins bool
	flag.BoolVar(&origins, "o", false, "list of allowed origins")
	flag.Parse()

	var originsList []string
	if origins {
		originsList = flag.Args()
	} else {
		originsList = []string{"http://localhost:3000"}
	}

	server.Serve(originsList)
}
