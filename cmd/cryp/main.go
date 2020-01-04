package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"../../src/server"
	"../../src/stream"

	_ "github.com/gorilla/handlers"
	_ "github.com/gorilla/mux"
	_ "github.com/gorilla/websocket"
)

func main() {
	var (
		c           = make(chan os.Signal, 1)
		ctx, cancel = context.WithCancel(context.Background())
	)

	signal.Notify(c, os.Interrupt)

	go func() {
		go stream.Stream(ctx)
		<-c
		cancel()
	}()

	if err := server.Serve(ctx); err != nil {
		log.Printf("failed to serve: %+v\n", err)
	}
}
