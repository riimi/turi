package main

import (
	"log"
	"os"
)

func main() {
	log.Printf("Starting the service... build time: %s, release: %s", BuildTime, Release)
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "1323"
	}

	server := NewGameServer(Port(port))
	server.Run()
}
