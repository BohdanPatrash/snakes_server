package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/BohdanPatrash/snakes_server/pkg/websocket"
)

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", wsEndpoint(pool))
}

func wsEndpoint(pool *websocket.Pool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		connection, err := websocket.Upgrade(w, r)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		client := &websocket.Client{
			ID:   strconv.Itoa(len(pool.Clients)),
			Conn: connection,
			Pool: pool,
		}

		pool.Register <- client
		client.Read()
	}
}

func main() {
	setupRoutes()
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
