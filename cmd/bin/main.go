package main

import (
	"bytes"
	"compress/flate"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"

	"github.com/gorilla/websocket"
)

func main() {
	var wg sync.WaitGroup
	var urls []url.URL

	subs := map[string][]byte{
		"real.okex.com:8443":      []byte(`{"op":"subscribe", "args": ["spot/depth5:ETH-USDT"]}`),
		"api.huobi.pro":           []byte(`{"sub": "market.btcusdt.depth.step0"}`),
		"stream.binance.com:9443": []byte(`{"method": "SUBSCRIBE", "params": ["ethbtc@depth"], "id": 1}`),
		"www.bitmex.com":          []byte(`{"op": "subscribe", "args": ["orderBookL2:ETHUSD"]}`)}

	urls = append(urls,
		url.URL{Scheme: "wss", Host: "real.okex.com:8443", Path: "/ws/v3", RawQuery: "compress=true"},
		url.URL{Scheme: "wss", Host: "api.huobi.pro", Path: "ws"},
		url.URL{Scheme: "wss", Host: "www.bitmex.com", Path: "realtime"},
		url.URL{Scheme: "wss", Host: "stream.binance.com:9443", Path: "/ws/ethbtc@depth20"})

	// for graceful shutdown
	shutdown := make(chan struct{})
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		<-interrupt
		log.Println("Interrupt.")
		close(shutdown)
	}()

	for _, u := range urls {
		wg.Add(1)
		go con(u, shutdown, subs[u.Host], &wg)
	}

	wg.Wait()
}

func con(u url.URL, shutdown chan struct{}, sub []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Printf("connecting to %s", u.Host)
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	c.WriteMessage(websocket.TextMessage, sub)
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			messageType, message, err := c.ReadMessage()
			switch messageType {
			case websocket.TextMessage:
				log.Printf("recv: %s", message)
			case websocket.BinaryMessage:
				text, err := gzip(message)
				if err != nil {
					log.Println("err", err)
				} else {
					log.Printf("recv: %s", text)
				}
			}
			if err != nil {
				log.Println("read:", err)
				return
			}
		}
	}()

	for {
		select {
		case <-done:
			return

		case <-shutdown:
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			return
		}
	}

}

func gzip(in []byte) ([]byte, error) {
	reader := flate.NewReader(bytes.NewReader(in))
	defer reader.Close()
	return ioutil.ReadAll(reader)
}
