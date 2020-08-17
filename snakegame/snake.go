package snakegame

import (
	"fmt"
	"log"

	guuid "github.com/google/uuid"
)

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Snake struct {
	PlayerID string  `json:"-"`
	HeadX    int     `json:"headX"`
	HeadY    int     `json:"headY"`
	Length   int     `json:"length"`
	SpeedX   int     `json:"speedX"`
	SpeedY   int     `json:"speedY"`
	Tail     []Point `json:"tail"`
	IsDead   bool    `json:"isDead"`
}

type snakeGeneration struct {
	x   int
	y   int
	spX int
	spY int
}

//generateSnakes generates hardcoded positions for 4 players only.
// This can be changed in the future for more flexible generation
func generateSnakes(cnf Config) map[int]snakeGeneration {
	snakes := map[int]snakeGeneration{}
	chunks := int((cnf.GamefieldSize - 1) / cnf.Scale)
	mid := int(chunks / 2)
	snakes[0] = snakeGeneration{
		spX: 1,
		spY: 0,
		x:   2 * cnf.Scale,
		y:   (mid + 1) * cnf.Scale,
	}
	snakes[1] = snakeGeneration{
		spX: 0,
		spY: 1,
		x:   (mid) * cnf.Scale,
		y:   2 * cnf.Scale,
	}
	snakes[2] = snakeGeneration{
		spX: -1,
		spY: 0,
		x:   (chunks - 2) * cnf.Scale,
		y:   (mid) * cnf.Scale,
	}
	snakes[3] = snakeGeneration{
		spX: 0,
		spY: -1,
		x:   (mid + 1) * cnf.Scale,
		y:   (chunks - 2) * cnf.Scale,
	}
	return snakes
}

//NewSnake creates new snake for pos(position) from 0 to 3.
// This function generates only 4 hardcoded positions this fuction might
// be changed in the future to beter generate first positions of snakes
func NewSnake(pos int, cnf Config) (*Snake, error) {
	id, err := guuid.NewRandom()
	if err != nil {
		log.Printf("could not create UUID: %v", err)
	}
	if pos > 3 || pos < 0 {
		return nil, fmt.Errorf("The position of the snake shoud be from 0 to 3")
	}
	posSnake := generateSnakes(cnf)[pos]
	tail := []Point{{}, {}, {}}
	for i := 2; i >= 0; i-- {
		tail[i].X = posSnake.x + posSnake.spX*cnf.Scale*i
		tail[i].Y = posSnake.y + posSnake.spY*cnf.Scale*i
	}
	snake := &Snake{
		PlayerID: id.String(),
		Length:   4,
		SpeedX:   posSnake.spX,
		SpeedY:   posSnake.spY,
		Tail:     tail,
		IsDead:   false,
		HeadX:    posSnake.x + posSnake.spX*cnf.Scale*3,
		HeadY:    posSnake.y + posSnake.spY*cnf.Scale*3,
	}
	return snake, nil
}

func (snake *Snake) Move(cnf Config) {
	snake.moveTail()
	snake.moveHead(cnf.Scale)
}

//CheckCollision checks if snake has collision for the current state of the snake and snakes passed as variable
func (snake *Snake) CheckCollision(cnf Config, snakes []Snake) {
	if snake.hitsWall(cnf) || snake.hitsSnakes(snakes) {
		snake.die()
	}
}

func (snake *Snake) moveHead(scl int) {
	snake.HeadX += snake.SpeedX * scl
	snake.HeadY += snake.SpeedY * scl
}

func (snake *Snake) moveTail() {
	tail := []Point{
		Point{
			X: snake.HeadX,
			Y: snake.HeadY,
		},
	}
	snake.Tail = append(tail, snake.Tail[:len(snake.Tail)-1]...)
}

func (snake *Snake) die() {
	snake.IsDead = true
}

func (snake *Snake) hitsWall(cnf Config) bool {
	if snake.HeadX < 0 ||
		snake.HeadY < 0 ||
		snake.HeadX > (cnf.CanvasWidth-cnf.Scale-1) ||
		snake.HeadY > (cnf.CanvasHeight-cnf.Scale-1) {
		return true
	}
	return false
}

func (snake *Snake) hitsSnakes(snakes []Snake) bool {
	for _, s := range snakes {
		if snake.HeadX == s.HeadX && snake.HeadY == s.HeadY && snake.PlayerID != s.PlayerID {
			return true
		}
		for _, tail := range s.Tail {
			if snake.HeadX == tail.X && snake.HeadY == tail.Y {
				return true
			}
		}
	}
	return false
}
