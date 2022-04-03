package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/mjekson/http_service/internals/app"
	"github.com/mjekson/http_service/internals/cfg"
)

func main() {
	config := cfg.LoadAndStoreConfig()

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	server := app.NewServer(config, ctx)

	go func() {
		osCall := <-c
		log.Printf("system call:%v", osCall)
		server.Shutdown()
		cancel()
	}()

	server.Serve()
}
