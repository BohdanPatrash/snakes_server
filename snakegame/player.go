package snakegame

import (
	"log"

	"github.com/gorilla/websocket"
)

type Player struct {
	ID          string
	Conn        *websocket.Conn
	GameSession *GameSession
	Snake       *Snake
}

func NewPlayer(session *GameSession, connection *websocket.Conn, snake *Snake) *Player {
	player := &Player{
		ID:          snake.PlayerID,
		GameSession: session,
		Conn:        connection,
		Snake:       snake,
	}
	return player
}

func (p *Player) Read() {
	defer func() {
		p.GameSession.Unregister <- p
		p.Conn.Close()
	}()

	for {
		message := DirectionInfo{}
		err := p.Conn.ReadJSON(&message)
		if err != nil {
			log.Println("could not read JSON: ", err)
			return
		}
		message.PlayerID = p.ID
		p.GameSession.UpdateSnake <- message
		// fmt.Printf("message received: %+v\n", message)
	}
}
