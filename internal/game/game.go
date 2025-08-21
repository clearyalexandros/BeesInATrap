package game

import "fmt"

type Game struct {
	Player *Player
	Hive   []*Bee
	Turns  int
}

// NewGame creates a new game with initialized player and bee hive
func NewGame() *Game {
	game := &Game{
		Player: NewPlayer(),
		Hive:   make([]*Bee, 0, 31), // 1 Queen + 5 Workers + 25 Drones
		Turns:  0,
	}

	game.initializeHive()
	return game
}

// initializeHive creates the bee hive according to specifications
func (g *Game) initializeHive() {
	// Add 1 Queen Bee
	g.Hive = append(g.Hive, NewBee(Queen))

	// Add 5 Worker Bees
	for i := 0; i < 5; i++ {
		g.Hive = append(g.Hive, NewBee(Worker))
	}

	// Add 25 Drone Bees
	for i := 0; i < 25; i++ {
		g.Hive = append(g.Hive, NewBee(Drone))
	}
}

// GetAliveBees returns a slice of all living bees
func (g *Game) GetAliveBees() []*Bee {
	var aliveBees []*Bee
	for _, bee := range g.Hive {
		if bee.IsAlive {
			aliveBees = append(aliveBees, bee)
		}
	}
	return aliveBees
}

// GetBeesByType returns a slice of all living bees of a specific type
func (g *Game) GetBeesByType(beeType BeeType) []*Bee {
	var bees []*Bee
	for _, bee := range g.Hive {
		if bee.IsAlive && bee.Type == beeType {
			bees = append(bees, bee)
		}
	}
	return bees
}

// IsGameOver checks if the game has ended
func (g *Game) IsGameOver() bool {
	// Player is dead
	if !g.Player.IsAlive() {
		return true
	}

	// All bees are dead
	aliveBees := g.GetAliveBees()
	return len(aliveBees) == 0
}

// KillAllBees kills all remaining bees (when Queen dies)
func (g *Game) KillAllBees() {
	for _, bee := range g.Hive {
		if bee.IsAlive {
			bee.IsAlive = false
			bee.HP = 0
		}
	}
}

// PrintGameStatus displays current game state
func (g *Game) PrintGameStatus() {
	fmt.Printf("\n=== Game Status ===\n")
	fmt.Printf("Player HP: %d/%d\n", g.Player.HP, g.Player.MaxHP)

	queens := g.GetBeesByType(Queen)
	workers := g.GetBeesByType(Worker)
	drones := g.GetBeesByType(Drone)

	fmt.Printf("Alive Bees:\n")
	fmt.Printf("  Queens: %d\n", len(queens))
	fmt.Printf("  Workers: %d\n", len(workers))
	fmt.Printf("  Drones: %d\n", len(drones))
	fmt.Printf("Turns: %d\n", g.Turns)
	fmt.Println("==================")
}

// Start begins the game
func (g *Game) Start() {
	fmt.Println("Welcome to Bees in the Trap!")
	fmt.Println("Your mission: Destroy the hive before the bees sting you to death!")
	fmt.Println("Type 'hit' to attack the hive, or 'auto' to let the game run automatically.")
	g.PrintGameStatus()
}
