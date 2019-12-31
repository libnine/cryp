package stream

import (
	"bytes"
	"compress/flate"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type binance struct {
	Host         string          `json:"host,omitempty"`
	LastUpdateID int             `json:"lastUpdateId"`
	Bids         [][]json.Number `json:"bids"`
	Asks         [][]json.Number `json:"asks"`
}

type bitmex struct {
	Host   string     `json:"host,omitempty"`
	Table  string     `json:"table"`
	Action string     `json:"action,omitempty"`
	Data   []bitmexl2 `json:"data"`
}

type bitmexl2 struct {
	Symbol string  `json:"symbol"`
	ID     int64   `json:"id"`
	Price  float64 `json:"price,omitempty"`
	Side   string  `json:"side"`
	Size   int64   `json:"size"`
}

type bitmexInit struct {
	ready       bool
	Table       string          `json:"table"`
	Action      string          `json:"action"`
	Keys        []string        `json:"keys"`
	Types       json.RawMessage `json:"types,omitempty"`
	ForeignKeys json.RawMessage `json:"foreignKeys,omitempty"`
	Attributes  json.RawMessage `json:"attributes,omitempty"`
	Filter      json.RawMessage `json:"filter,omitempty"`
	Data        json.RawMessage `json:"data"`
}

type okex struct {
	Host  string   `json:"host,omitempty"`
	Table string   `json:"table"`
	Data  []okexl2 `json:"data"`
}

type okexl2 struct {
	Asks [][]json.Number `json:"asks"`
	Bids [][]json.Number `json:"bids"`
	Inst string          `json:"instrument_id"`
	Ts   *time.Time      `json:"timestamp"`
}

// global variables
var (
	BitmexTable     = bitmexInit{}
	IncomingOkex    = make(chan okex)
	IncomingBinance = make(chan binance)
	IncomingBitmex  = make(chan bitmex)
)

// Stream for crypto data
func Stream(ctx context.Context) {
	var (
		shutdown = make(chan struct{})
		urls     []url.URL
		wg       sync.WaitGroup
	)

	subs := map[string][]byte{
		// "api.huobi.pro":           []byte(`{"sub": "market.btcusdt.depth.step0"}`),
		"real.okex.com:8443":      []byte(`{"op":"subscribe", "args": ["spot/depth5:ETH-USDT"]}`),
		"stream.binance.com:9443": []byte(`{"method": "SUBSCRIBE", "params": ["ethusdt@depth"], "id": 1}`),
		"www.bitmex.com":          []byte(`{"op": "subscribe", "args": ["orderBookL2_25:ETHUSD"]}`)}

	urls = append(urls,
		// url.URL{Scheme: "wss", Host: "api.huobi.pro", Path: "ws"},
		url.URL{Scheme: "wss", Host: "real.okex.com:8443", Path: "/ws/v3", RawQuery: "compress=true"},
		url.URL{Scheme: "wss", Host: "www.bitmex.com", Path: "realtime?subscribe=instrument"},
		url.URL{Scheme: "wss", Host: "stream.binance.com:9443", Path: "/ws/ethusdt@depth20"})

	go func() {
		<-ctx.Done()
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

	log.Printf("connecting to: %+s\n", u.Host)
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	c.WriteMessage(websocket.TextMessage, sub)

	var (
		attempts int
		done     = make(chan struct{})
	)

	go func() {
		defer close(done)
		for {
			messageType, message, err := c.ReadMessage()
			switch messageType {
			case websocket.TextMessage:
				switch u.Host {

				case "real.okex.com:8443":
					var x okex
					_ = json.Unmarshal(message, &x)
					x.Host = "okex"
					IncomingOkex <- x

				case "stream.binance.com:9443":
					var x binance
					_ = json.Unmarshal(message, &x)
					x.Host = "binance"
					IncomingBinance <- x

				case "www.bitmex.com":
					if !BitmexTable.ready {
						attempts++
						_ = json.Unmarshal(message, &BitmexTable)

						if len(BitmexTable.Keys) > 0 {
							BitmexTable.ready = true
						}

						if attempts >= 10 {
							log.Printf("stopped bitmex after 10 attempts\n")
							return
						}

					} else {
						var x bitmex
						_ = json.Unmarshal(message, &x)
						x.Host = "bitmex"
						IncomingBitmex <- x
					}
				}

			case websocket.BinaryMessage:
				reader := flate.NewReader(bytes.NewReader(message))
				defer reader.Close()
				text, err := ioutil.ReadAll(reader)
				if err != nil {
					log.Printf("err: %s %s", u.Host, err)
				}

				if u.Host == "real.okex.com:8443" {
					var x okex
					_ = json.Unmarshal(text, &x)
					x.Host = "okex"
					IncomingOkex <- x
				}
			}

			if err != nil {
				log.Printf("err: %+v\n", err)
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
				log.Printf("write close:+%v\n", err)
				return
			}
			return
		}
	}
}
