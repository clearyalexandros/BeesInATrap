package game

// Bee configuration constants
const (
	// Queen Bee stats
	QueenHP          = 100
	QueenDamage      = 10
	QueenTakesDamage = 10

	// Worker Bee stats
	WorkerHP          = 75
	WorkerDamage      = 5
	WorkerTakesDamage = 25

	// Drone Bee stats
	DroneHP          = 60
	DroneDamage      = 1
	DroneTakesDamage = 30
)

type BeeType int

const (
	Queen BeeType = iota
	Worker
	Drone
)

// BeeStats holds all the stats for a particular bee type
type BeeStats struct {
	HP          int
	Damage      int
	TakesDamage int
}

// BeeStatsTable provides O(1) lookup for bee stats by type (map access vs switch statements)
var BeeStatsTable = map[BeeType]BeeStats{
	Queen:  {HP: QueenHP, Damage: QueenDamage, TakesDamage: QueenTakesDamage},
	Worker: {HP: WorkerHP, Damage: WorkerDamage, TakesDamage: WorkerTakesDamage},
	Drone:  {HP: DroneHP, Damage: DroneDamage, TakesDamage: DroneTakesDamage},
}

type Bee struct {
	Type   BeeType
	HP     int
	MaxHP  int
	Damage int
}

// NewBee creates a new bee with stats based on what type it is
func NewBee(beeType BeeType) *Bee {
	stats := BeeStatsTable[beeType]
	return &Bee{
		Type:   beeType,
		HP:     stats.HP,
		MaxHP:  stats.HP,
		Damage: stats.Damage,
	}
}

// IsAlive checks if the bee still has health left
func (b *Bee) IsAlive() bool {
	return b.HP > 0
}

// TakeDamage hits the bee and deals damage based on what type it is
func (b *Bee) TakeDamage() {
	stats := BeeStatsTable[b.Type]
	b.HP -= stats.TakesDamage
	if b.HP < 0 {
		b.HP = 0
	}
}

// String returns the name of the bee type as a string
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
