package main

import (
	"flag"
	"log"

	client "github.com/chesnokovilya/faraway/client/lib"
)

func main() {
	serverAddr := flag.String("server", "localhost:8080", "Server address")
	flag.Parse()
	powClient := client.NewClient(*serverAddr)
	if err := powClient.Connect(); err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
}
