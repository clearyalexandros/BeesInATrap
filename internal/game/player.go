package game

type Player struct {
	HP    int
	MaxHP int
}

// NewPlayer creates a new player with 100 HP
func NewPlayer() *Player {
	return &Player{
		HP:    100,
		MaxHP: 100,
	}
}

// TakeDamage reduces the player's HP by the specified amount
func (p *Player) TakeDamage(damage int) {
	p.HP -= damage
	if p.HP < 0 {
		p.HP = 0
	}
}

// IsAlive returns true if the player has HP remaining
func (p *Player) IsAlive() bool {
	return p.HP > 0
}
