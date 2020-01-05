package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	exch "github.com/libnine/cryp/internal/exchange"
)

func main() {
	var (
		c           = make(chan os.Signal, 1)
		ctx, cancel = context.WithCancel(context.Background())
	)

	signal.Notify(c, os.Interrupt)

	go func() {
		go exch.Feed(ctx)
		<-c
		cancel()
	}()

	if err := exch.Serve(ctx); err != nil {
		log.Printf("failed to serve: %+v\n", err)
	}
}
