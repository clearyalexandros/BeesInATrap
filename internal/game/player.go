package game

// Player configuration constants
const (
	PlayerStartingHP = 100
)

type Player struct {
	HP    int
	MaxHP int
}

// NewPlayer creates a new player starting with full health
func NewPlayer() Player {
	return Player{
		HP:    PlayerStartingHP,
		MaxHP: PlayerStartingHP,
	}
}

// TakeDamage hurts the player and reduces their health
func (p *Player) TakeDamage(damage int) {
	p.HP -= damage
	if p.HP < 0 {
		p.HP = 0
	}
}

// IsAlive checks if the player still has health left
func (p Player) IsAlive() bool {
	return p.HP > 0
}
