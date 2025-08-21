package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Game struct {
	Player   *Player
	Hive     []*Bee
	Turns    int
	AutoMode bool
	rng      *rand.Rand
}

// NewGame sets up a fresh game with a player and a full hive of bees
func NewGame() *Game {
	game := &Game{
		Player:   NewPlayer(),
		Hive:     make([]*Bee, 0, 31), // Space for 1 Queen + 5 Workers + 25 Drones
		Turns:    0,
		AutoMode: false,
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	game.initializeHive()
	return game
} // initializeHive populates the hive with all the bees according to the game rules
func (g *Game) initializeHive() {
	// Add the one and only Queen Bee
	g.Hive = append(g.Hive, NewBee(Queen))

	// Add the Worker Bees
	for i := 0; i < 5; i++ {
		g.Hive = append(g.Hive, NewBee(Worker))
	}

	// Add the Drone Bees
	for i := 0; i < 25; i++ {
		g.Hive = append(g.Hive, NewBee(Drone))
	}
}

// GetAliveBees gives you all the bees that are still buzzing around
func (g *Game) GetAliveBees() []*Bee {
	var aliveBees []*Bee
	for _, bee := range g.Hive {
		if bee.IsAlive {
			aliveBees = append(aliveBees, bee)
		}
	}
	return aliveBees
}

// GetBeesByType finds all living bees of a particular type
func (g *Game) GetBeesByType(beeType BeeType) []*Bee {
	var bees []*Bee
	for _, bee := range g.Hive {
		if bee.IsAlive && bee.Type == beeType {
			bees = append(bees, bee)
		}
	}
	return bees
}

// IsGameOver checks if someone has won or lost the game
func (g *Game) IsGameOver() bool {
	// Player is dead
	if !g.Player.IsAlive() {
		return true
	}

	// All bees are dead
	aliveBees := g.GetAliveBees()
	return len(aliveBees) == 0
}

// KillAllBees wipes out the entire hive (happens when the Queen dies)
func (g *Game) KillAllBees() {
	for _, bee := range g.Hive {
		if bee.IsAlive {
			bee.IsAlive = false
			bee.HP = 0
		}
	}
}

// PrintGameStatus shows the current state of the battle
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

// Start welcomes the player and shows them what's happening
func (g *Game) Start() {
	fmt.Println("Welcome to Bees in the Trap!")
	fmt.Println("Your mission: Destroy the hive before the bees sting you to death!")
	fmt.Println("Type 'hit' to attack the hive, or 'auto' to let the game run automatically.")
	g.PrintGameStatus()
}

// PlayGame keeps the game running until someone wins or loses
func (g *Game) PlayGame() {
	scanner := bufio.NewScanner(os.Stdin)

	for !g.IsGameOver() {
		if g.AutoMode {
			// Let the computer play automatically
			g.PlayerTurn("hit")
			time.Sleep(500 * time.Millisecond) // Small pause so you can follow along
		} else {
			// Wait for the player to tell us what to do
			fmt.Print("\nEnter command (hit/auto/quit): ")
			if !scanner.Scan() {
				break
			}

			input := strings.TrimSpace(strings.ToLower(scanner.Text()))

			switch input {
			case "hit":
				g.PlayerTurn(input)
			case "auto":
				fmt.Println("Switching to auto mode...")
				g.AutoMode = true
				continue
			case "quit":
				fmt.Println("Thanks for playing!")
				return
			default:
				fmt.Println("Invalid command. Use 'hit', 'auto', or 'quit'.")
				continue
			}
		}

		// See if the game ended after the player's turn
		if g.IsGameOver() {
			break
		}

		// Now it's the bees' turn to fight back
		g.BeeTurn()
	}

	g.EndGame()
}

// PlayerTurn lets the player do something on their turn
func (g *Game) PlayerTurn(command string) {
	g.Turns++
	fmt.Printf("\n--- Turn %d: Player Turn ---\n", g.Turns)

	if command == "hit" {
		g.PlayerAttack()
	}
}

// PlayerAttack makes the player swing at the hive
func (g *Game) PlayerAttack() {
	aliveBees := g.GetAliveBees()
	if len(aliveBees) == 0 {
		fmt.Println("No bees left to attack!")
		return
	}

	// Sometimes you miss completely
	if g.rng.Float64() < 0.15 {
		fmt.Println("Miss! You just missed the hive, better luck next time!")
		return
	}

	// Pick a random bee to hit
	targetBee := aliveBees[g.rng.Intn(len(aliveBees))]

	fmt.Printf("Direct Hit! You attacked a %s bee!\n", targetBee.Type.String())

	// Hit the bee
	targetBee.TakeDamage()

	if !targetBee.IsAlive {
		fmt.Printf("You killed the %s bee! (%d damage dealt)\n", targetBee.Type.String(), g.getDamageDealtTo(targetBee.Type))

		// Special rule: killing the Queen kills everyone
		if targetBee.Type == Queen {
			fmt.Println("🔥 QUEEN BEE ELIMINATED! All remaining bees flee in terror! 🔥")
			g.KillAllBees()
		}
	} else {
		fmt.Printf("The %s bee took %d damage and has %d HP remaining.\n", targetBee.Type.String(), g.getDamageDealtTo(targetBee.Type), targetBee.HP)
	}
}

// BeeTurn makes the bees attack back
func (g *Game) BeeTurn() {
	fmt.Printf("\n--- Turn %d: Bees Turn ---\n", g.Turns)

	aliveBees := g.GetAliveBees()
	if len(aliveBees) == 0 {
		return
	}

	// Sometimes the bees all miss you
	if g.rng.Float64() < 0.20 {
		randomBee := aliveBees[g.rng.Intn(len(aliveBees))]
		fmt.Printf("Buzz! That was close! The %s Bee just missed you!\n", randomBee.Type.String())
		return
	}

	// Pick a random bee to sting you
	attackingBee := aliveBees[g.rng.Intn(len(aliveBees))]

	fmt.Printf("Sting! You just got stung by a %s bee!\n", attackingBee.Type.String())

	// Take damage from the bee
	g.Player.TakeDamage(attackingBee.Damage)

	fmt.Printf("You took %d damage and now have %d HP remaining.\n", attackingBee.Damage, g.Player.HP)

	if !g.Player.IsAlive() {
		fmt.Println("💀 You have been stung to death! 💀")
	}
}

// getDamageDealtTo tells you how much damage each bee type takes when hit
func (g *Game) getDamageDealtTo(beeType BeeType) int {
	switch beeType {
	case Queen:
		return 10
	case Worker:
		return 25
	case Drone:
		return 30
	default:
		return 0
	}
}

// EndGame shows the final results and says goodbye
func (g *Game) EndGame() {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("                 GAME OVER")
	fmt.Println(strings.Repeat("=", 50))

	if g.Player.IsAlive() {
		fmt.Println("🎉 CONGRATULATIONS! YOU WON! 🎉")
		fmt.Printf("You successfully destroyed the hive in %d turns!\n", g.Turns)
	} else {
		fmt.Println("💀 GAME OVER - YOU DIED 💀")
		fmt.Printf("The bees defeated you after %d turns.\n", g.Turns)
	}

	// Show how the battle went
	fmt.Println("\n--- GAME SUMMARY ---")
	fmt.Printf("Total turns: %d\n", g.Turns)
	fmt.Printf("Final player HP: %d/%d\n", g.Player.HP, g.Player.MaxHP)

	aliveBees := g.GetAliveBees()
	fmt.Printf("Bees remaining: %d/31\n", len(aliveBees))

	if len(aliveBees) > 0 {
		queens := g.GetBeesByType(Queen)
		workers := g.GetBeesByType(Worker)
		drones := g.GetBeesByType(Drone)
		fmt.Printf("  Queens: %d, Workers: %d, Drones: %d\n", len(queens), len(workers), len(drones))
	}

	fmt.Println("\nThanks for playing Bees in the Trap!")
}
