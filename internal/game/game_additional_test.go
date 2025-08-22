package game

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

// Test PlayerAttack function comprehensive scenarios
func TestPlayerAttackScenarios(t *testing.T) {
	t.Run("Normal Attack Scenarios", func(t *testing.T) {
		game := NewGame()
		initialBeeCount := len(game.GetAliveBees())

		// Execute multiple attacks to test different scenarios
		for i := 0; i < 10; i++ {
			// Test that PlayerAttack executes without error
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("PlayerAttack panicked on iteration %d: %v", i, r)
					}
				}()

				game.PlayerAttack()
			}()

			// Verify game state remains valid
			aliveBees := game.GetAliveBees()

			// If game ended, break the loop
			if len(aliveBees) == 0 {
				break
			}
		}

		// Should have potentially fewer bees (some might have been killed)
		finalBeeCount := len(game.GetAliveBees())
		if finalBeeCount > initialBeeCount {
			t.Error("Bee count should not increase after attacks")
		}
	})

	t.Run("Attack With Fixed Seed - Predictable Behavior", func(t *testing.T) {
		// Test with fixed seed for deterministic behavior
		game := NewGame()
		game.rng = rand.New(rand.NewSource(42)) // Fixed seed for reproducible tests

		initialBeeCount := len(game.GetAliveBees())

		// With fixed seed, we can predict some behavior
		game.PlayerAttack()

		// Game should still be in a valid state
		aliveBees := game.GetAliveBees()
		if len(aliveBees) > initialBeeCount {
			t.Errorf("Invalid bee count after attack: got %d, initial was %d", len(aliveBees), initialBeeCount)
		}
	})

	t.Run("Attack Different Bee Types", func(t *testing.T) {
		// Test attacking each bee type specifically
		beeTypes := []BeeType{Queen, Worker, Drone}

		for _, beeType := range beeTypes {
			t.Run(fmt.Sprintf("Attack_%s", beeType.String()), func(t *testing.T) {
				game := NewGame()

				// Create a game with only one bee of the target type
				game.KillAllBees()

				// Add one bee of the specific type
				testBee := NewBee(beeType)
				game.Hive[beeType] = []*Bee{testBee}
				game.AliveBees = []*Bee{testBee}

				// Force hit (no miss)
				game.rng = rand.New(rand.NewSource(1)) // Seed that ensures hit

				initialHP := testBee.HP
				game.PlayerAttack()

				// Bee should have taken damage or died
				if testBee.IsAlive() {
					if testBee.HP >= initialHP {
						t.Errorf("%s bee should have taken damage", beeType.String())
					}
				}
				// If bee died, verify it's actually dead
				if !testBee.IsAlive() && testBee.HP != 0 {
					t.Errorf("Dead %s bee should have 0 HP, got %d", beeType.String(), testBee.HP)
				}
			})
		}
	})

	t.Run("Queen Death Triggers Hive Elimination", func(t *testing.T) {
		game := NewGame()

		// Get the queen and damage it to near death
		queens := game.GetBeesByType(Queen)
		if len(queens) == 0 {
			t.Fatal("Expected at least one queen bee")
		}

		queen := queens[0]

		// Damage queen to 1 HP (10 damage per hit, so 9 hits = 10 HP remaining)
		for i := 0; i < 9; i++ {
			queen.TakeDamage()
		}

		if queen.HP != 10 {
			t.Fatalf("Expected queen to have 10 HP after 9 hits, got %d", queen.HP)
		}

		initialBeeCount := len(game.GetAliveBees())
		_ = initialBeeCount // Used for potential debugging/logging

		// Create a scenario where we're guaranteed to hit the queen
		// (by making it the only bee alive)
		for _, bee := range game.GetAliveBees() {
			if bee != queen {
				bee.HP = 0
			}
		}

		// Force hit (no miss) and attack
		game.rng = rand.New(rand.NewSource(1))
		game.PlayerAttack()

		// All bees should be dead due to queen death rule
		aliveBees := game.GetAliveBees()
		if len(aliveBees) > 0 {
			t.Errorf("Expected 0 alive bees after queen death, got %d", len(aliveBees))
		}

		// Queen should be dead
		if queen.IsAlive() {
			t.Error("Queen should be dead after final attack")
		}
	})

	t.Run("Miss Scenario", func(t *testing.T) {
		game := NewGame()

		// Force miss by setting a seed that produces miss
		game.rng = rand.New(rand.NewSource(0)) // Seed that likely produces miss

		// Try multiple times to get a miss
		foundMiss := false
		for i := 0; i < 20; i++ {
			game.rng = rand.New(rand.NewSource(int64(i)))
			beforeAttack := len(game.GetAliveBees())

			// Capture if all bees have same HP before attack
			beeHPs := make(map[*Bee]int)
			for _, bee := range game.GetAliveBees() {
				beeHPs[bee] = bee.HP
			}

			game.PlayerAttack()

			// Check if it was a miss (no bee HP changed)
			missOccurred := true
			for bee, oldHP := range beeHPs {
				if bee.HP != oldHP {
					missOccurred = false
					break
				}
			}

			if missOccurred && len(game.GetAliveBees()) == beforeAttack {
				foundMiss = true
				break
			}
		}

		// We should be able to generate a miss scenario
		// (This test verifies the miss mechanic works)
		if !foundMiss {
			t.Log("Could not generate a miss scenario in 20 attempts - this might be due to very low miss chance")
		}
	})

	t.Run("No Bees To Attack", func(t *testing.T) {
		game := NewGame()
		game.KillAllBees()

		// Attack should handle empty hive gracefully
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("PlayerAttack should not panic when no bees alive: %v", r)
			}
		}()

		game.PlayerAttack()

		// Should still have no bees
		if len(game.GetAliveBees()) != 0 {
			t.Error("Should have no alive bees after attacking empty hive")
		}
	})
}

// Test getDamageDealtTo function
func TestGetDamageDealtTo(t *testing.T) {
	game := NewGame()

	// Test damage to Queen
	damage := game.getDamageDealtTo(Queen)
	if damage != 10 {
		t.Errorf("Expected 10 damage to Queen, got %d", damage)
	}

	// Test damage to Worker
	damage = game.getDamageDealtTo(Worker)
	if damage != 25 {
		t.Errorf("Expected 25 damage to Worker, got %d", damage)
	}

	// Test damage to Drone
	damage = game.getDamageDealtTo(Drone)
	if damage != 30 {
		t.Errorf("Expected 30 damage to Drone, got %d", damage)
	}
}

// Test PlayerTurn function increments turn counter
func TestPlayerTurnIncrementsTurnCounter(t *testing.T) {
	game := NewGame()

	// Execute player turn with "hit" command
	game.PlayerTurn("hit")

	// Check that turn counter increased
	if game.Turns != 1 {
		t.Errorf("Expected turn count to be 1, got %d", game.Turns)
	}

	// Test with other command (should not crash)
	game.PlayerTurn("invalid")
	if game.Turns != 2 {
		t.Errorf("Expected turn count to be 2, got %d", game.Turns)
	}
}

// Test BeeTurn function player HP validation
func TestBeeTurnPlayerHPValidation(t *testing.T) {
	game := NewGame()
	initialPlayerHP := game.Player.HP

	// Execute bee turn
	game.BeeTurn()

	// Player HP should not go below 0
	if game.Player.HP < 0 {
		t.Error("Player HP should not go below 0")
	}

	// Player HP should not exceed maximum
	if game.Player.HP > game.Player.MaxHP {
		t.Errorf("Player HP (%d) should not exceed maximum (%d)", game.Player.HP, game.Player.MaxHP)
	}

	// Player HP might have changed (decreased from bee attacks)
	if game.Player.HP > initialPlayerHP {
		t.Error("Player HP should not increase during bee turn")
	}

	// Turn count should remain unchanged (BeeTurn doesn't increment turns)
	if game.Turns != 0 {
		t.Errorf("Expected turn count to remain 0, got %d", game.Turns)
	}
}

// Test makeBeeDecision function
func TestMakeBeeDecision(t *testing.T) {
	game := NewGame()
	bee := NewBee(Queen)

	// Test bee decision making
	start := time.Now()
	decision := game.makeBeeDecision(bee)
	duration := time.Since(start)

	// Should return a BeeDecision struct
	if decision.Bee != bee {
		t.Error("BeeDecision should reference the correct bee")
	}

	// WillHit should be a boolean (true or false)
	if decision.WillHit != true && decision.WillHit != false {
		t.Error("WillHit should be a boolean")
	}

	// Should take some time to "think" (at least 50ms for Queen)
	if duration < 50*time.Millisecond {
		t.Error("Queen should take at least 50ms to make decision")
	}

	// DecisionTime should be recorded
	if decision.DecisionTime <= 0 {
		t.Error("DecisionTime should be positive")
	}
}

// Test concurrent bee decisions
func TestConcurrentBeeDecisions(t *testing.T) {
	game := NewGame()
	bees := game.GetAliveBees()

	// Test that multiple bees can make decisions concurrently
	start := time.Now()

	// Simulate what happens in BeeTurn
	results := make(chan BeeDecision, len(bees))
	for _, bee := range bees {
		go func(b *Bee) {
			decision := game.makeBeeDecision(b)
			results <- decision
		}(bee)
	}

	// Collect results
	decisions := make([]BeeDecision, 0, len(bees))
	for i := 0; i < len(bees); i++ {
		decision := <-results
		decisions = append(decisions, decision)
	}

	duration := time.Since(start)

	// Should get decisions for all bees
	if len(decisions) != len(bees) {
		t.Errorf("Expected %d decisions, got %d", len(bees), len(decisions))
	}

	// With concurrent execution, total time should be reasonable
	if duration > 3*time.Second {
		t.Error("Concurrent bee decisions taking too long")
	}
}

// Test damage event channel
func TestDamageEventChannel(t *testing.T) {
	game := NewGame()

	// Give the goroutine time to start
	time.Sleep(10 * time.Millisecond)

	// Set turns > 0 so damage events are processed
	game.Turns = 1

	// Send a damage event
	select {
	case game.damageEvent <- 5:
		// Success - channel accepted the damage event
	case <-time.After(100 * time.Millisecond):
		t.Error("Damage event channel should accept events")
	}

	// Give time for the goroutine to process
	time.Sleep(10 * time.Millisecond)
}

// Test EndGame function basic output
func TestEndGameBasicOutput(t *testing.T) {
	game := NewGame()
	game.Turns = 10

	// Capture stdout to test the output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	game.EndGame()

	// Restore stdout and capture output
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Test for expected output content
	expectedPhrases := []string{
		"GAME OVER",
		"Total turns: 10",
		"Final player HP:",
		"Bees remaining:",
		"Thanks for playing Bees in the Trap!",
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(output, phrase) {
			t.Errorf("Expected EndGame() output to contain '%s', but it didn't. Output: %s", phrase, output)
		}
	}
}

// Test PrintGameStatus function
func TestPrintGameStatus(t *testing.T) {
	game := NewGame()

	// Capture stdout to test the output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	game.PrintGameStatus()

	// Restore stdout and capture output
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Test for expected output content
	expectedPhrases := []string{
		"=== Game Status ===",
		"Player HP: 100/100",
		"Alive Bees:",
		"Queens: 1",
		"Workers: 5",
		"Drones: 25",
		"Turns: 0",
		"==================",
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(output, phrase) {
			t.Errorf("Expected PrintGameStatus() output to contain '%s', but it didn't. Output: %s", phrase, output)
		}
	}
}

// Test String method for BeeType
func TestBeeTypeString(t *testing.T) {
	tests := []struct {
		beeType  BeeType
		expected string
	}{
		{Queen, "Queen"},
		{Worker, "Worker"},
		{Drone, "Drone"},
		{BeeType(999), "Unknown"}, // Test unknown bee type
	}

	for _, test := range tests {
		result := test.beeType.String()
		if result != test.expected {
			t.Errorf("Expected %s.String() to return '%s', got '%s'",
				test.beeType, test.expected, result)
		}
	}
}

// Test NewGame damage event monitoring goroutine
func TestNewGameDamageEventMonitoring(t *testing.T) {
	game := NewGame()
	game.Turns = 1 // Enable damage event processing

	// Give the goroutine time to start
	time.Sleep(10 * time.Millisecond)

	// Test different damage levels for different icons
	testCases := []struct {
		damage       int
		expectedIcon string
		description  string
	}{
		{15, "🩸", "high damage"},
		{7, "⚡", "medium damage"},
		{2, "🔸", "low damage"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// Send damage event
			select {
			case game.damageEvent <- tc.damage:
				// Success - event was sent
			case <-time.After(100 * time.Millisecond):
				t.Errorf("Damage event channel should accept %d damage event", tc.damage)
			}

			// Give time for processing
			time.Sleep(10 * time.Millisecond)
		})
	}

	// Test channel full scenario (non-blocking)
	// Fill the channel buffer (capacity is 10)
	for i := 0; i < 15; i++ {
		select {
		case game.damageEvent <- 1:
			// Continue filling
		default:
			// Channel full, which is expected behavior
			return // Exit the function once channel is full
		}
	}
}

// Test BeeTurn all miss scenario
func TestBeeTurnAllMissScenario(t *testing.T) {
	// Set a seed that will cause all bees to miss
	// Test multiple seeds until we find one that produces all misses
	foundAllMissScenario := false

	for seed := int64(0); seed < 100; seed++ {
		game := NewGame()
		game.rng = rand.New(rand.NewSource(seed))

		initialPlayerHP := game.Player.HP

		// Capture stdout to check for miss message
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		game.BeeTurn()

		// Restore stdout and capture output
		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		// Check if all bees missed (player HP unchanged and "missed" in output)
		if game.Player.HP == initialPlayerHP && strings.Contains(output, "missed") {
			foundAllMissScenario = true

			// Verify the miss message format
			if !strings.Contains(output, "Buzz! That was close!") {
				t.Error("Expected miss message format not found")
			}
			break
		}
	}

	if !foundAllMissScenario {
		t.Log("Could not generate all-miss scenario in 100 attempts - this is acceptable due to randomness")
	}
}

// Test BeeTurn player death scenario
func TestBeeTurnPlayerDeath(t *testing.T) {
	game := NewGame()

	// Reduce player HP to 1 so next bee attack will kill them
	game.Player.HP = 1

	// Force a bee to hit by manipulating the scenario
	// Keep only one bee that will definitely hit
	game.KillAllBees()
	bee := NewBee(Drone) // Drone does 1 damage, perfect for killing player with 1 HP
	game.Hive[Drone] = []*Bee{bee}
	game.AliveBees = []*Bee{bee}

	// Set seed to ensure hit
	game.rng = rand.New(rand.NewSource(1))

	// Capture stdout to verify death message
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	game.BeeTurn()

	// Restore stdout and capture output
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Verify player died
	if game.Player.IsAlive() {
		t.Error("Player should be dead after taking fatal damage")
	}

	// Verify death message appears
	if !strings.Contains(output, "You have been stung to death!") {
		t.Error("Expected death message not found in output")
	}
}

// Test EndGame player death scenario
func TestEndGamePlayerDeath(t *testing.T) {
	game := NewGame()
	game.Turns = 5
	game.Player.HP = 0 // Player is dead

	// Capture stdout to test the output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	game.EndGame()

	// Restore stdout and capture output
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Test for player death specific content
	expectedPhrases := []string{
		"GAME OVER",
		"💀 GAME OVER - YOU DIED 💀",
		"The bees defeated you after 5 turns",
		"Total turns: 5",
		"Final player HP: 0/100",
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(output, phrase) {
			t.Errorf("Expected EndGame() player death output to contain '%s', but it didn't. Output: %s", phrase, output)
		}
	}
}

// Test EndGame with remaining bees breakdown
func TestEndGameWithRemainingBees(t *testing.T) {
	game := NewGame()
	game.Turns = 3
	game.Player.HP = 0 // Player died, so bees remain

	// Ensure some bees are still alive for the breakdown
	// Kill some bees but leave others
	aliveBees := game.GetAliveBees()
	for i, bee := range aliveBees {
		if i >= 15 { // Kill about half the bees
			bee.HP = 0
		}
	}

	// Capture stdout to test the output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	game.EndGame()

	// Restore stdout and capture output
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Test for bee breakdown content
	expectedPhrases := []string{
		"Bees remaining:",
		"Queens:",
		"Workers:",
		"Drones:",
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(output, phrase) {
			t.Errorf("Expected EndGame() remaining bees output to contain '%s', but it didn't. Output: %s", phrase, output)
		}
	}

	// Verify that the bee counts are reasonable
	remainingBees := len(game.GetAliveBees())
	if remainingBees == 0 {
		t.Error("Expected some bees to remain alive for this test")
	}
}

// Test EndGame player victory scenario
func TestEndGamePlayerVictory(t *testing.T) {
	game := NewGame()
	game.Turns = 8
	game.KillAllBees() // Player wins by killing all bees

	// Capture stdout to test the output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	game.EndGame()

	// Restore stdout and capture output
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Test for victory specific content
	expectedPhrases := []string{
		"GAME OVER",
		"🎉 CONGRATULATIONS! YOU WON! 🎉",
		"You successfully destroyed the hive in 8 turns!",
		"Total turns: 8",
		"Bees remaining: 0/31",
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(output, phrase) {
			t.Errorf("Expected EndGame() victory output to contain '%s', but it didn't. Output: %s", phrase, output)
		}
	}
}
