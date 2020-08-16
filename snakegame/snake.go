package snakegame

type Point struct {
	X int
	Y int
}

type Snake struct {
	ClientID string  `json:"-"`
	HeadX    int     `json:"headX,omitempty"`
	HeadY    int     `json:"headY,omitempty"`
	Length   int     `json:"length,omitempty"`
	SpeedX   int     `json:"speedX,omitempty"`
	SpeedY   int     `json:"speedY,omitempty"`
	Tail     []Point `json:"tail,omitempty"`
	IsDead   bool    `json:"isDead,omitempty"`
}

func (snake *Snake) Move(cnf Config) {
	snake.moveHead(cnf.Scale)
	snake.moveTail()
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
	for _, enemySnake := range snakes {
		if snake.HeadX == enemySnake.HeadX && snake.HeadY == enemySnake.HeadY {
			return true
		}
		for _, tail := range enemySnake.Tail {
			if snake.HeadX == tail.X && snake.HeadY == tail.Y {
				return true
			}
		}
	}
	return false
}
