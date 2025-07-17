package main

import (
	"log"

	"github.com/sternerr/termtalk/internal/server"
)

func main() {
	srv, err := server.NewServer("0.0.0.0", "6969")
	if err != nil {
		log.Fatal(err.Error())
	}
	srv.Listen()
}
