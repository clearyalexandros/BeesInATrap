package game

import "testing"

func TestNewPlayer(t *testing.T) {
	player := NewPlayer()

	if player.HP != 100 {
		t.Errorf("Expected player HP to be 100, got %d", player.HP)
	}

	if player.MaxHP != 100 {
		t.Errorf("Expected player MaxHP to be 100, got %d", player.MaxHP)
	}

	if !player.IsAlive() {
		t.Error("Expected new player to be alive")
	}
}

func TestNewGame(t *testing.T) {
	game := NewGame()

	// Test player initialization - Player is now a value type, not pointer
	if game.Player.HP != PlayerStartingHP {
		t.Errorf("Expected player to start with %d HP, got %d", PlayerStartingHP, game.Player.HP)
	}

	// Test hive initialization - now using total alive bees count
	aliveBees := game.GetAliveBees()
	if len(aliveBees) != TotalBees {
		t.Errorf("Expected hive to have %d bees, got %d", TotalBees, len(aliveBees))
	}

	// Count bee types
	queens := game.GetBeesByType(Queen)
	workers := game.GetBeesByType(Worker)
	drones := game.GetBeesByType(Drone)

	if len(queens) != 1 {
		t.Errorf("Expected 1 Queen bee, got %d", len(queens))
	}

	if len(workers) != 5 {
		t.Errorf("Expected 5 Worker bees, got %d", len(workers))
	}

	if len(drones) != 25 {
		t.Errorf("Expected 25 Drone bees, got %d", len(drones))
	}
}

func TestBeeStats(t *testing.T) {
	// Test Queen bee stats
	queen := NewBee(Queen)
	if queen.HP != 100 || queen.MaxHP != 100 || queen.Damage != 10 {
		t.Errorf("Queen stats incorrect: HP=%d, MaxHP=%d, Damage=%d", queen.HP, queen.MaxHP, queen.Damage)
	}

	// Test Worker bee stats
	worker := NewBee(Worker)
	if worker.HP != 75 || worker.MaxHP != 75 || worker.Damage != 5 {
		t.Errorf("Worker stats incorrect: HP=%d, MaxHP=%d, Damage=%d", worker.HP, worker.MaxHP, worker.Damage)
	}

	// Test Drone bee stats
	drone := NewBee(Drone)
	if drone.HP != 60 || drone.MaxHP != 60 || drone.Damage != 1 {
		t.Errorf("Drone stats incorrect: HP=%d, MaxHP=%d, Damage=%d", drone.HP, drone.MaxHP, drone.Damage)
	}
}

func TestBeeTakeDamage(t *testing.T) {
	// Test Queen taking damage
	queen := NewBee(Queen)
	queen.TakeDamage()
	if queen.HP != 90 {
		t.Errorf("Expected Queen to have 90 HP after taking damage, got %d", queen.HP)
	}

	// Test Worker taking damage
	worker := NewBee(Worker)
	worker.TakeDamage()
	if worker.HP != 50 {
		t.Errorf("Expected Worker to have 50 HP after taking damage, got %d", worker.HP)
	}

	// Test Drone taking damage
	drone := NewBee(Drone)
	drone.TakeDamage()
	if drone.HP != 30 {
		t.Errorf("Expected Drone to have 30 HP after taking damage, got %d", drone.HP)
	}
}

func TestQueenBeeDamage(t *testing.T) {
	queen := NewBee(Queen)

	// Test initial state
	if queen.HP != 100 || queen.MaxHP != 100 {
		t.Errorf("Queen should start with 100/100 HP, got %d/%d", queen.HP, queen.MaxHP)
	}

	if !queen.IsAlive() {
		t.Error("Queen should be alive initially")
	}

	// Test taking damage multiple times (Queen takes 10 damage per hit)
	for i := 1; i <= 9; i++ {
		queen.TakeDamage()
		expectedHP := 100 - (i * 10)
		if queen.HP != expectedHP {
			t.Errorf("After %d hits, Queen should have %d HP, got %d", i, expectedHP, queen.HP)
		}
		if !queen.IsAlive() {
			t.Errorf("Queen should still be alive after %d hits", i)
		}
	}

	// Final hit should kill the Queen
	queen.TakeDamage()
	if queen.HP != 0 {
		t.Errorf("Queen should have 0 HP after 10 hits, got %d", queen.HP)
	}
	if queen.IsAlive() {
		t.Error("Queen should be dead after 10 hits")
	}
}

func TestWorkerBeeDamage(t *testing.T) {
	worker := NewBee(Worker)

	// Test initial state
	if worker.HP != 75 || worker.MaxHP != 75 {
		t.Errorf("Worker should start with 75/75 HP, got %d/%d", worker.HP, worker.MaxHP)
	}

	if !worker.IsAlive() {
		t.Error("Worker should be alive initially")
	}

	// First hit (Worker takes 25 damage per hit)
	worker.TakeDamage()
	if worker.HP != 50 {
		t.Errorf("After 1 hit, Worker should have 50 HP, got %d", worker.HP)
	}
	if !worker.IsAlive() {
		t.Error("Worker should still be alive after 1 hit")
	}

	// Second hit
	worker.TakeDamage()
	if worker.HP != 25 {
		t.Errorf("After 2 hits, Worker should have 25 HP, got %d", worker.HP)
	}
	if !worker.IsAlive() {
		t.Error("Worker should still be alive after 2 hits")
	}

	// Third hit should kill the Worker
	worker.TakeDamage()
	if worker.HP != 0 {
		t.Errorf("Worker should have 0 HP after 3 hits, got %d", worker.HP)
	}
	if worker.IsAlive() {
		t.Error("Worker should be dead after 3 hits")
	}
}

func TestDroneBeeDamage(t *testing.T) {
	drone := NewBee(Drone)

	// Test initial state
	if drone.HP != 60 || drone.MaxHP != 60 {
		t.Errorf("Drone should start with 60/60 HP, got %d/%d", drone.HP, drone.MaxHP)
	}

	if !drone.IsAlive() {
		t.Error("Drone should be alive initially")
	}

	// First hit (Drone takes 30 damage per hit)
	drone.TakeDamage()
	if drone.HP != 30 {
		t.Errorf("After 1 hit, Drone should have 30 HP, got %d", drone.HP)
	}
	if !drone.IsAlive() {
		t.Error("Drone should still be alive after 1 hit")
	}

	// Second hit should kill the Drone
	drone.TakeDamage()
	if drone.HP != 0 {
		t.Errorf("Drone should have 0 HP after 2 hits, got %d", drone.HP)
	}
	if drone.IsAlive() {
		t.Error("Drone should be dead after 2 hits")
	}
}

func TestBeeTypeDamageValues(t *testing.T) {
	tests := []struct {
		beeType        BeeType
		expectedHP     int
		expectedMaxHP  int
		expectedDamage int
		damagePerHit   int
		hitsToKill     int
	}{
		{Queen, 100, 100, 10, 10, 10},
		{Worker, 75, 75, 5, 25, 3},
		{Drone, 60, 60, 1, 30, 2},
	}

	for _, test := range tests {
		t.Run(test.beeType.String(), func(t *testing.T) {
			bee := NewBee(test.beeType)

			// Test initial stats
			if bee.HP != test.expectedHP {
				t.Errorf("%s should have %d HP, got %d", test.beeType.String(), test.expectedHP, bee.HP)
			}
			if bee.MaxHP != test.expectedMaxHP {
				t.Errorf("%s should have %d MaxHP, got %d", test.beeType.String(), test.expectedMaxHP, bee.MaxHP)
			}
			if bee.Damage != test.expectedDamage {
				t.Errorf("%s should deal %d damage, got %d", test.beeType.String(), test.expectedDamage, bee.Damage)
			}

			// Test damage progression
			for hit := 1; hit < test.hitsToKill; hit++ {
				bee.TakeDamage()
				expectedHP := test.expectedHP - (hit * test.damagePerHit)
				if bee.HP != expectedHP {
					t.Errorf("After %d hits, %s should have %d HP, got %d", hit, test.beeType.String(), expectedHP, bee.HP)
				}
				if !bee.IsAlive() {
					t.Errorf("%s should still be alive after %d hits", test.beeType.String(), hit)
				}
			}

			// Final hit should kill
			bee.TakeDamage()
			if bee.HP != 0 {
				t.Errorf("%s should have 0 HP after %d hits, got %d", test.beeType.String(), test.hitsToKill, bee.HP)
			}
			if bee.IsAlive() {
				t.Errorf("%s should be dead after %d hits", test.beeType.String(), test.hitsToKill)
			}
		})
	}
}

func TestBeeExcessiveDamage(t *testing.T) {
	// Test that HP doesn't go below 0 when taking excessive damage
	tests := []struct {
		beeType BeeType
		name    string
	}{
		{Queen, "Queen"},
		{Worker, "Worker"},
		{Drone, "Drone"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bee := NewBee(test.beeType)

			// Kill the bee multiple times
			for i := 0; i < 20; i++ {
				bee.TakeDamage()
			}

			if bee.HP != 0 {
				t.Errorf("%s HP should not go below 0, got %d", test.name, bee.HP)
			}
			if bee.IsAlive() {
				t.Errorf("%s should be dead after excessive damage", test.name)
			}
		})
	}
}

func TestPlayerTakeDamage(t *testing.T) {
	player := NewPlayer()

	// Test normal damage
	player.TakeDamage(25)
	if player.HP != 75 {
		t.Errorf("Expected player to have 75 HP after taking 25 damage, got %d", player.HP)
	}

	// Test fatal damage
	player.TakeDamage(100)
	if player.HP != 0 {
		t.Errorf("Expected player HP to be 0 after fatal damage, got %d", player.HP)
	}

	if player.IsAlive() {
		t.Error("Expected player to be dead after fatal damage")
	}
}

func TestGameOver(t *testing.T) {
	game := NewGame()

	// Game should not be over at start
	if game.IsGameOver() {
		t.Error("Game should not be over at start")
	}

	// Kill player
	game.Player.TakeDamage(100)
	if !game.IsGameOver() {
		t.Error("Game should be over when player is dead")
	}

	// Reset and kill all bees
	game2 := NewGame()
	game2.KillAllBees()
	if !game2.IsGameOver() {
		t.Error("Game should be over when all bees are dead")
	}
}

func TestQueenDeathRule(t *testing.T) {
	game := NewGame()

	// Find and kill the queen
	queens := game.GetBeesByType(Queen)
	if len(queens) != 1 {
		t.Fatalf("Expected 1 queen, got %d", len(queens))
	}

	queen := queens[0]
	// Kill queen (takes 10 hits of 10 damage each)
	for i := 0; i < 10; i++ {
		queen.TakeDamage()
	}

	if queen.IsAlive() {
		t.Error("Queen should be dead after taking full damage")
	}

	// Manually trigger the queen death rule
	game.KillAllBees()

	// Check all bees are dead
	aliveBees := game.GetAliveBees()
	if len(aliveBees) != 0 {
		t.Errorf("Expected 0 alive bees after queen death, got %d", len(aliveBees))
	}
}
