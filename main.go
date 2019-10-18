package main

import (
	"log"
	"os"
)

func main() {
	log.Printf("Starting the service... build time: %s, commit: %s, release: %s", BuildTime, Commit, Release)
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "1323"
	}

	server := NewGameServer(Port(port))
	server.Run()
}
