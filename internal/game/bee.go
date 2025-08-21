package game

type BeeType int

const (
	Queen  BeeType = 0
	Worker BeeType = 1
	Drone  BeeType = 2
)

type Bee struct {
	Type    BeeType
	HP      int
	MaxHP   int
	Damage  int
	IsAlive bool
}

// NewBee creates a new bee with the appropriate stats based on type
func NewBee(beeType BeeType) *Bee {
	bee := &Bee{
		Type:    beeType,
		IsAlive: true,
	}

	switch beeType {
	case Queen:
		bee.HP = 100
		bee.MaxHP = 100
		bee.Damage = 10
	case Worker:
		bee.HP = 75
		bee.MaxHP = 75
		bee.Damage = 5
	case Drone:
		bee.HP = 60
		bee.MaxHP = 60
		bee.Damage = 1
	}

	return bee
}

// TakeDamage applies damage to the bee based on its type
func (b *Bee) TakeDamage() {
	var damage int
	switch b.Type {
	case Queen:
		damage = 10
	case Worker:
		damage = 25
	case Drone:
		damage = 30
	}

	b.HP -= damage
	if b.HP <= 0 {
		b.HP = 0
		b.IsAlive = false
	}
}

// String returns a string representation of the bee type
func (bt BeeType) String() string {
	switch bt {
	case Queen:
		return "Queen"
	case Worker:
		return "Worker"
	case Drone:
		return "Drone"
	default:
		return "Unknown"
	}
}
