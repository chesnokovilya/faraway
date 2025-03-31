package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chesnokovilya/faraway/cmd/server"
)

func main() {
	tcpServer, err := server.NewServer(":8080")
	if err != nil {
		log.Fatal(err)
	}
	if err := tcpServer.Start(); err != nil {
		log.Fatal(err)
	}
	defer tcpServer.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
