package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	client "github.com/chesnokovilya/faraway/client/lib"
)

func main() {
	serverAddr := flag.String("server", "localhost:8080", "Server address")
	flag.Parse()

	powClient := client.NewClient(*serverAddr)

	start := time.Now()
	if err := powClient.Connect(); err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	elapsed := time.Since(start)

	fmt.Printf("Completed in %v\n", elapsed)
}
