package snakegame

import (
	"fmt"
	"log"
	"time"
)

type Config struct {
	CanvasWidth   int `json:"canvasWidth"`
	CanvasHeight  int `json:"canvasHeight"`
	GamefieldSize int `json:"gamefieldSize"`
	Scale         int `json:"scale"`
}

func DefaultConfig() *Config {
	return &Config{
		CanvasWidth:   600,
		CanvasHeight:  440,
		GamefieldSize: 401,
		Scale:         20,
	}
}

type GameSession struct {
	Register    chan *Player
	Unregister  chan *Player
	Players     map[string]*Player
	UpdateSnake chan DirectionInfo
	Config      *Config
}

func NewGameSession(cnf *Config) *GameSession {
	return &GameSession{
		Register:    make(chan *Player),
		Unregister:  make(chan *Player),
		Players:     make(map[string]*Player),
		UpdateSnake: make(chan DirectionInfo),
		Config:      cnf,
	}
}

func (session *GameSession) Start() {
	ticker := time.NewTicker(125 * time.Millisecond)
	for {
		select {
		//registers new client
		case player := <-session.Register:
			session.Players[player.ID] = player
			fmt.Printf("Player %s: has joined session\n", player.ID)
			// for client := range pool.Clients {
			// 	//change broadcasting message
			// 	// client.Conn.WriteJSON()
			// }
		//unregisters client
		case player := <-session.Unregister:
			delete(session.Players, player.ID)
			fmt.Printf("Player %s: has disconnected from session\n", player.ID)
			// for client := range pool.Clients {
			// 	//change broadcasting message
			// 	// client.Conn.WriteJSON()
			// }
		case direction := <-session.UpdateSnake:
			UpdatePlayerDirection(session, direction)
		case <-ticker.C:
			SendTickUpdate(session)
		}
	}
}

func UpdatePlayerDirection(session *GameSession, direction DirectionInfo) {
	snake := session.Players[direction.PlayerID].Snake
	// fmt.Println()
	// fmt.Println(snake)
	// fmt.Println()
	snake.SpeedX = direction.SpeedX
	snake.SpeedY = direction.SpeedY
}

func SendTickUpdate(session *GameSession) {
	updateSnakes(session)
	sendData(session)
}

func updateSnakes(session *GameSession) {
	snakesToCheck := []*Snake{}
	snakesAlive := []Snake{}
	for _, player := range session.Players {
		if player.Snake.IsDead {
			continue
		}
		player.Snake.Move(*session.Config)
		snakesToCheck = append(snakesToCheck, player.Snake)
		snakesAlive = append(snakesAlive, *player.Snake)
	}

	for _, snake := range snakesToCheck {
		snake.CheckCollision(*session.Config, snakesAlive)
	}
}

func sendData(session *GameSession) {
	for _, player := range session.Players {
		data := OutboundData{
			Snake:       player.Snake,
			EnemySnakes: []Snake{},
		}
		for _, p := range session.Players {
			if player.ID == p.ID {
				continue
			}
			data.EnemySnakes = append(data.EnemySnakes, *p.Snake)
		}
		err := player.Conn.WriteJSON(data)
		if err != nil {
			log.Println(err)
		}
	}
}
