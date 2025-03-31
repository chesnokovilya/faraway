package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chesnokovilya/faraway/cmd/server"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		cancel()
	}()
	g, gCtx := errgroup.WithContext(ctx)
	tcpServer, err := server.NewServer(":8080")
	if err != nil {
		log.Fatal(err)
	}
	g.Go(func() error {
		if err := tcpServer.Start(); err != nil {
			log.Fatal(err)
		}
		return nil
	})
	g.Go(func() error {
		<-gCtx.Done()
		log.Println("server is going to shutdown")
		tcpServer.Stop()
		return nil
	})
	if err := g.Wait(); err != nil {
		fmt.Printf("exit reason: %s \n", err)
	}
}
