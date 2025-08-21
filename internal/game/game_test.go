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

	// Test player initialization
	if game.Player == nil {
		t.Error("Expected game to have a player")
	}

	if game.Player.HP != 100 {
		t.Errorf("Expected player to start with 100 HP, got %d", game.Player.HP)
	}

	// Test hive initialization
	if len(game.Hive) != 31 {
		t.Errorf("Expected hive to have 31 bees, got %d", len(game.Hive))
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
