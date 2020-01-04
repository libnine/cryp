package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	stream "../stream"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	clients = make(map[*websocket.Conn]bool)
	handler = handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(r)
	r       = mux.NewRouter()
	srv     = &http.Server{
		Addr:         ":8000",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// Serve ws data coming from crypto exchanges
func Serve(ctx context.Context) (err error) {
	r.HandleFunc("/ids", idsHandler).Methods("GET")
	r.HandleFunc("/ws", wsHandler).Methods("GET")

	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %+s\n", err)
		}
	}()

	log.Printf("server started")

	go func() {
		echo(ctx)
	}()

	<-ctx.Done()

	log.Printf("server stopped")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("server shutdown failed: %+s", err)
	}

	log.Printf("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}

	return
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

func echo(ctx context.Context) {
	for {
		select {
		case v := <-stream.IncomingOkex:
			for client := range clients {
				w, err := client.NextWriter(websocket.TextMessage)
				if err != nil {
					client.Close()
					delete(clients, client)
					continue
				}

				err = json.NewEncoder(w).Encode(&v)
			}

		case v := <-stream.IncomingBinance:
			for client := range clients {
				w, err := client.NextWriter(websocket.TextMessage)
				if err != nil {
					client.Close()
					delete(clients, client)
					continue
				}

				err = json.NewEncoder(w).Encode(&v)
			}

		case v := <-stream.IncomingBitmex:
			for client := range clients {
				w, err := client.NextWriter(websocket.TextMessage)
				if err != nil {
					client.Close()
					delete(clients, client)
					continue
				}

				err = json.NewEncoder(w).Encode(&v)
			}

		case v := <-stream.IncomingBitstamp:
			for client := range clients {
				w, err := client.NextWriter(websocket.TextMessage)
				if err != nil {
					client.Close()
					delete(clients, client)
					continue
				}

				err = json.NewEncoder(w).Encode(&v)
			}

		case <-ctx.Done():
			for client := range clients {
				client.Close()
			}

			return
		}
	}
}
