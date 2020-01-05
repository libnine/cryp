package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	btc "github.com/libnine/cryp/internal/pkg/cryp"
)

func main() {
	var (
		c           = make(chan os.Signal, 1)
		ctx, cancel = context.WithCancel(context.Background())
	)

	signal.Notify(c, os.Interrupt)

	go func() {
		go btc.Feed(ctx)
		<-c
		cancel()
	}()

	if err := btc.Serve(ctx); err != nil {
		log.Printf("failed to serve: %+v\n", err)
	}
}
