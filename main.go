package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/BohdanPatrash/snakes_server/snakegame"

	"github.com/BohdanPatrash/snakes_server/websocket"
	guuid "github.com/google/uuid"
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
		}
		id, err := guuid.NewRandom()
		if err != nil {
			log.Printf("could not create UUID: %v", err)
		}

		snake := &snakegame.Snake{
			ClientID: id.String(),
			SpeedX:   0,
			SpeedY:   0,
		}

		player := &snakegame.Player{
			ID:          id.String(),
			Conn:        connection,
			GameSession: session,
			Snake:       snake,
		}

		session.Register <- player
		player.Read()
	}
}

func main() {
	setupRoutes()
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
