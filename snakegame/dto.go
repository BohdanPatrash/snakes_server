package snakegame

type HandshakeResponse struct {
	Color  string `json:"color,omitempty"`
	Snake  Snake  `json:"snake"`
	Config `json:"config"`
}

type DirectionInfo struct {
	PlayerID string `json:"-"`
	SpeedX   int    `json:"speedX,omitempty"`
	SpeedY   int    `json:"speedY,omitempty"`
}

type OutboundData struct {
	Snake       *Snake  `json:"snake"`
	EnemySnakes []Snake `json:"enemySnakes"`
}
