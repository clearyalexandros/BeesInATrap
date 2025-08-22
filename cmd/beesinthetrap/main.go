package main

import (
	"flag"
	"fmt"

	"github.com/clearyalexandros/BeesInATrap/internal/game"
)

func main() {
	// Define command-line flags
	playerHP := flag.Int("player-hp", 100, "Starting health points for the player")
	playerMissChance := flag.Float64("player-miss", 0.15, "Player miss chance (0.0-1.0)")
	beesMissChance := flag.Float64("bees-miss", 0.20, "Bees miss chance (0.0-1.0)")
	autoDelay := flag.Int("auto-delay", 500, "Auto mode delay in milliseconds")

	// Hive composition flags
	queenCount := flag.Int("queens", 1, "Number of Queen bees in the hive")
	workerCount := flag.Int("workers", 5, "Number of Worker bees in the hive")
	droneCount := flag.Int("drones", 25, "Number of Drone bees in the hive")

	// Help flag
	showHelp := flag.Bool("help", false, "Show help information")

	flag.Parse()

	if *showHelp {
		fmt.Println("🐝 Bees in the Trap - Configuration Options")
		fmt.Println("==========================================")
		flag.PrintDefaults()
		fmt.Println("\nExample usage:")
		fmt.Println("  beesinthetrap --player-hp 150 --player-miss 0.10 --bees-miss 0.30")
		fmt.Println("  beesinthetrap --queens 2 --workers 10 --drones 50")
		fmt.Println("  beesinthetrap --auto-delay 1000 --help")
		return
	}

	// Validate input ranges
	if *playerHP <= 0 {
		fmt.Println("Error: Player HP must be greater than 0")
		return
	}
	if *playerMissChance < 0.0 || *playerMissChance > 1.0 {
		fmt.Println("Error: Player miss chance must be between 0.0 and 1.0")
		return
	}
	if *beesMissChance < 0.0 || *beesMissChance > 1.0 {
		fmt.Println("Error: Bees miss chance must be between 0.0 and 1.0")
		return
	}
	if *autoDelay < 0 {
		fmt.Println("Error: Auto delay must be non-negative")
		return
	}
	if *queenCount < 0 || *workerCount < 0 || *droneCount < 0 {
		fmt.Println("Error: Bee counts must be non-negative")
		return
	}

	fmt.Println("Starting Bees in the Trap...")

	// Create game configuration
	config := game.GameConfig{
		PlayerHP:         *playerHP,
		PlayerMissChance: *playerMissChance,
		BeesMissChance:   *beesMissChance,
		AutoModeDelay:    *autoDelay,
		QueenCount:       *queenCount,
		WorkerCount:      *workerCount,
		DroneCount:       *droneCount,
	}

	// Show configuration if any non-default values are used
	if *playerHP != 100 || *playerMissChance != 0.15 || *beesMissChance != 0.20 ||
		*autoDelay != 500 || *queenCount != 1 || *workerCount != 5 || *droneCount != 25 {
		fmt.Printf("Custom Configuration:\n")
		fmt.Printf("  Player HP: %d\n", *playerHP)
		fmt.Printf("  Player Miss Chance: %.1f%%\n", *playerMissChance*100)
		fmt.Printf("  Bees Miss Chance: %.1f%%\n", *beesMissChance*100)
		fmt.Printf("  Auto Mode Delay: %dms\n", *autoDelay)
		fmt.Printf("  Hive: %d Queens, %d Workers, %d Drones (%d total)\n",
			*queenCount, *workerCount, *droneCount, *queenCount+*workerCount+*droneCount)
		fmt.Println()
	}

	g := game.NewGameWithConfig(config)
	g.Start()

	// Let's play!
	g.PlayGame()
}
