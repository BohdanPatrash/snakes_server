package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/BohdanPatrash/snakes_server/snakegame"

	"github.com/BohdanPatrash/snakes_server/websocket"
)

func setupRoutes() {
	config := snakegame.DefaultConfig()
	session := snakegame.NewGameSession(config)
	go session.Start()

	http.HandleFunc("/ws", wsEndpoint(session))
}

func wsEndpoint(session *snakegame.GameSession) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		connection, err := websocket.Upgrade(w, r)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		snake, err := snakegame.NewSnake(len(session.Players), *session.Config)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		player := snakegame.NewPlayer(session, connection, snake)

		session.Register <- player
		player.Read()
	}
}

func main() {
	setupRoutes()
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
