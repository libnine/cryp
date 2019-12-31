package main

import (
	"sync"

	"../../src/server"
	"../../src/stream"
)

var (
	wg sync.WaitGroup
)

func main() {
	wg.Add(1)
	defer wg.Done()

	go stream.Stream(&wg)
	go server.Serve()

	wg.Wait()
}
