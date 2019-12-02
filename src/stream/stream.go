package stream

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Depth of all combined exchanges data
type Depth struct {
	depthBinance interface{}
	depthBitmex  interface{}
	depthOkex    interface{}
}

type binance struct {
	LastUpdateID int32       `json:"lastUpdateId"`
	Bids         [][]float64 `json:"bids"`
	Asks         [][]float64 `json:"asks"`
}

type bitmex struct {
	Table  string     `json:"table"`
	Action string     `json:"action"`
	Data   []bitmexl2 `json:"data"`
}

type bitmexl2 struct {
	Symbol string `json:"symbol"`
	ID     int64  `json:"id"`
	Side   string `json:"side"`
	Size   int64  `json:"size"`
}

type okex struct {
	Table string   `json:"table"`
	Data  []okexl2 `json:"data"`
}

type okexl2 struct {
	Asks [][]json.Number `json:"asks"`
	Bids [][]json.Number `json:"bids"`
	Inst string          `json:"instrument_id"`
	Ts   *time.Time      `json:"timestamp"`
}

func Stream() {
	var l2 Depth
	var wg sync.WaitGroup
	var urls []url.URL

	subs := map[string][]byte{
		"real.okex.com:8443":      []byte(`{"op":"subscribe", "args": ["spot/depth5:ETH-USDT"]}`),
		"api.huobi.pro":           []byte(`{"sub": "market.btcusdt.depth.step0"}`),
		"stream.binance.com:9443": []byte(`{"method": "SUBSCRIBE", "params": ["ethusdt@depth"], "id": 1}`),
		"www.bitmex.com":          []byte(`{"op": "subscribe", "args": ["orderBookL2:ETHUSD"]}`)}

	urls = append(urls,
		// url.URL{Scheme: "wss", Host: "real.okex.com:8443", Path: "/ws/v3", RawQuery: "compress=true"},
		// url.URL{Scheme: "wss", Host: "api.huobi.pro", Path: "ws"},
		// url.URL{Scheme: "wss", Host: "www.bitmex.com", Path: "realtime"},
		url.URL{Scheme: "wss", Host: "stream.binance.com:9443", Path: "/ws/ethusdt@depth20"})

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
		go con(u, shutdown, subs[u.Host], &l2, &wg)
	}

	wg.Wait()
}

func con(u url.URL, shutdown chan struct{}, sub []byte, l *Depth, wg *sync.WaitGroup) {
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
				if u.Host == "real.okex.com:8443" {
					var x okex
					_ = json.Unmarshal(message, &x)
					// l.depthOkex = x
					log.Println(x)
				} else if u.Host == "stream.binance.com:9443" {
					var x binance
					_ = json.Unmarshal(message, &x)
					l.depthBinance = x
					log.Printf("%v", l)
				} else if u.Host == "www.bitmex.com" {
					var x bitmex
					_ = json.Unmarshal(message, &x)
					// l.depthOkex = x
					log.Println(x)
				}
			case websocket.BinaryMessage:
				text, err := gzip(message)
				if err != nil {
					log.Printf("err: %s %s", u.Host, err)
				}

				if u.Host == "real.okex.com:8443" {
					var x okex
					_ = json.Unmarshal(text, &x)
					l.depthOkex = x
					log.Println(l)
				}
			}
			if err != nil {
				log.Println("err:", err)
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
