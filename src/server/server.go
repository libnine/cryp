package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	stream "../stream"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	clients  = make(map[*websocket.Conn]bool)
	done     = make(chan os.Signal, 1)
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// Serve ws data coming from crypto exchanges
func Serve(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	r := mux.NewRouter()
	r.HandleFunc("/ids", idsHandler).Methods("GET")
	r.HandleFunc("/ws", wsHandler).Methods("GET")

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go echo()
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(r)))

	wg.Wait()
}

func idsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(stream.BitmexTable.Data)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	clients[ws] = true
}

func echo() {
	defer close(done)
	for {
		select {
		case v := <-stream.IncomingOkex:
			for client := range clients {
				w, err := client.NextWriter(websocket.TextMessage)
				if err != nil {
					client.Close()
					delete(clients, client)
					break
				}

				err = json.NewEncoder(w).Encode(&v)
			}

		case v := <-stream.IncomingBinance:
			for client := range clients {
				w, err := client.NextWriter(websocket.TextMessage)
				if err != nil {
					client.Close()
					delete(clients, client)
					break
				}
				err = json.NewEncoder(w).Encode(&v)
			}

		case v := <-stream.IncomingBitmex:
			for client := range clients {
				w, err := client.NextWriter(websocket.TextMessage)
				if err != nil {
					client.Close()
					delete(clients, client)
					break
				}
				err = json.NewEncoder(w).Encode(&v)
			}

		case <-done:
			return
		}
	}
}
