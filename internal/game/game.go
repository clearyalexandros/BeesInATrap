package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

// Game configuration constants
const (
	PlayerMissChance = 0.15 // 15% chance for player to miss
	BeesMissChance   = 0.20 // 20% chance for all bees to miss
	AutoModeDelay    = 500  // Milliseconds to pause in auto mode

	// Hive composition
	QueenCount  = 1
	WorkerCount = 5
	DroneCount  = 25
	TotalBees   = QueenCount + WorkerCount + DroneCount
)

// BeeDecision represents a bee's decision to attack or miss
type BeeDecision struct {
	Bee          *Bee
	WillHit      bool
	DecisionTime time.Duration // How long the bee took to decide
}

type Game struct {
	Player      *Player            // Use pointer so we can modify the player
	Hive        map[BeeType][]*Bee // Map structure enables O(1) access to bees by type
	AliveBees   []*Bee             // Cached slice avoids O(n) scanning on each access
	Turns       int
	AutoMode    bool
	rng         *rand.Rand
	damageEvent chan int // Channel to signal damage events for stats monitoring
}

// NewGame sets up a fresh game with a player and a full hive of bees
func NewGame() *Game {
	game := &Game{
		Player:      &Player{HP: PlayerStartingHP, MaxHP: PlayerStartingHP},
		Hive:        make(map[BeeType][]*Bee),
		AliveBees:   make([]*Bee, 0, TotalBees),
		Turns:       0,
		AutoMode:    false,
		rng:         rand.New(rand.NewSource(time.Now().UnixNano())),
		damageEvent: make(chan int, 10), // Buffered channel for damage events
	}

	game.initializeHive()

	// Start event-driven game stats monitor
	go func() {
		for damage := range game.damageEvent {
			if game.Turns > 0 { // Only show stats after game starts
				aliveBees := len(game.GetAliveBees())
				survivalRate := float64(game.Player.HP) / float64(game.Player.MaxHP) * 100

				// Show different messages based on damage severity
				var damageIcon string
				switch {
				case damage >= 10:
					damageIcon = "🩸" // High damage
				case damage >= 5:
					damageIcon = "⚡" // Medium damage
				default:
					damageIcon = "🔸" // Low damage
				}

				fmt.Printf("%s Damage Alert: -%d HP | Turn %d | Player: %d/%d (%.1f%%) | Bees: %d\n",
					damageIcon, damage, game.Turns, game.Player.HP, game.Player.MaxHP, survivalRate, aliveBees)
			}
		}
	}()

	return game
} // initializeHive populates the hive with all the bees according to the game rules
func (g *Game) initializeHive() {
	// Initialize the map slices
	g.Hive[Queen] = make([]*Bee, 0, QueenCount)
	g.Hive[Worker] = make([]*Bee, 0, WorkerCount)
	g.Hive[Drone] = make([]*Bee, 0, DroneCount)

	// Add the Queen Bees
	for i := 0; i < QueenCount; i++ {
		bee := NewBee(Queen)
		g.Hive[Queen] = append(g.Hive[Queen], bee)
		g.AliveBees = append(g.AliveBees, bee)
	}

	// Add the Worker Bees
	for i := 0; i < WorkerCount; i++ {
		bee := NewBee(Worker)
		g.Hive[Worker] = append(g.Hive[Worker], bee)
		g.AliveBees = append(g.AliveBees, bee)
	}

	// Add the Drone Bees
	for i := 0; i < DroneCount; i++ {
		bee := NewBee(Drone)
		g.Hive[Drone] = append(g.Hive[Drone], bee)
		g.AliveBees = append(g.AliveBees, bee)
	}
}

// GetAliveBees gives you all the bees that are still alive
func (g *Game) GetAliveBees() []*Bee {
	// Rebuild the alive list by filtering out dead bees
	aliveBees := make([]*Bee, 0, len(g.AliveBees))
	for _, bee := range g.AliveBees {
		if bee.IsAlive() {
			aliveBees = append(aliveBees, bee)
		}
	}
	g.AliveBees = aliveBees // Update the cached list
	return aliveBees
}

// GetBeesByType finds all living bees of a particular type (O(1) map access to type group)
func (g *Game) GetBeesByType(beeType BeeType) []*Bee {
	var bees []*Bee
	for _, bee := range g.Hive[beeType] {
		if bee.IsAlive() {
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
	for _, beeList := range g.Hive {
		for _, bee := range beeList {
			if bee.IsAlive() {
				bee.HP = 0
			}
		}
	}
	g.AliveBees = []*Bee{} // Clear the alive list
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
			time.Sleep(AutoModeDelay * time.Millisecond) // Small pause so you can follow along
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
	if g.rng.Float64() < PlayerMissChance {
		fmt.Println("Miss! You just missed the hive, better luck next time!")
		return
	}

	// Pick a random bee to hit
	targetBee := aliveBees[g.rng.Intn(len(aliveBees))]

	fmt.Printf("Direct Hit! You attacked a %s bee!\n", targetBee.Type.String())

	// Hit the bee
	targetBee.TakeDamage()

	if !targetBee.IsAlive() {
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

// BeeTurn makes the bees attack back using concurrent decision making
func (g *Game) BeeTurn() {
	fmt.Printf("\n--- Turn %d: Bees Turn ---\n", g.Turns)

	aliveBees := g.GetAliveBees()
	if len(aliveBees) == 0 {
		return
	}

	// Channel to collect bee decisions
	decisionChan := make(chan BeeDecision, len(aliveBees))
	var wg sync.WaitGroup

	// Each bee makes a decision concurrently
	for _, bee := range aliveBees {
		wg.Add(1)
		go func(b *Bee) {
			defer wg.Done()
			decision := g.makeBeeDecision(b)
			decisionChan <- decision
		}(bee)
	}

	// Wait for all bees to make decisions
	go func() {
		wg.Wait()
		close(decisionChan)
	}()

	// Collect all decisions
	var hits []BeeDecision
	var misses []BeeDecision
	totalDecisionTime := time.Duration(0)

	for decision := range decisionChan {
		totalDecisionTime += decision.DecisionTime
		if decision.WillHit {
			hits = append(hits, decision)
		} else {
			misses = append(misses, decision)
		}
	}

	// Display thinking time (for demonstration)
	fmt.Printf("🧠 Bees consulted for %v total...\n", totalDecisionTime)

	// Execute attack based on decisions
	if len(hits) > 0 {
		// Random successful attack from the hits
		chosenAttack := hits[g.rng.Intn(len(hits))]
		fmt.Printf("Sting! You just got stung by a %s bee!\n", chosenAttack.Bee.Type.String())

		damage := chosenAttack.Bee.Damage
		g.Player.TakeDamage(damage)
		fmt.Printf("You took %d damage and now have %d HP remaining.\n", damage, g.Player.HP)

		// Trigger damage event for stats monitoring
		select {
		case g.damageEvent <- damage:
		default:
			// Channel full, skip this event (non-blocking)
		}

		if !g.Player.IsAlive() {
			fmt.Println("💀 You have been stung to death! 💀")
		}
	} else if len(misses) > 0 {
		// All bees missed - show a random miss
		chosenMiss := misses[g.rng.Intn(len(misses))]
		fmt.Printf("Buzz! That was close! The %s Bee just missed you!\n",
			chosenMiss.Bee.Type.String())
	}
}

// makeBeeDecision simulates a bee making an attack decision concurrently
func (g *Game) makeBeeDecision(bee *Bee) BeeDecision {
	start := time.Now()

	// Simulate different thinking times based on bee type
	var thinkingTime time.Duration
	switch bee.Type {
	case Queen:
		thinkingTime = time.Duration(50+g.rng.Intn(100)) * time.Millisecond // 50-150ms
	case Worker:
		thinkingTime = time.Duration(20+g.rng.Intn(60)) * time.Millisecond // 20-80ms
	case Drone:
		thinkingTime = time.Duration(10+g.rng.Intn(40)) * time.Millisecond // 10-50ms
	}

	// Simulate thinking
	time.Sleep(thinkingTime)

	// Make the hit/miss decision
	willHit := g.rng.Float64() >= BeesMissChance

	return BeeDecision{
		Bee:          bee,
		WillHit:      willHit,
		DecisionTime: time.Since(start),
	}
}

// getDamageDealtTo tells you how much damage each bee type takes when hit
func (g *Game) getDamageDealtTo(beeType BeeType) int {
	return BeeStatsTable[beeType].TakesDamage
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
