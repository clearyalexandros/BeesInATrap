package game

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// Test PlayerAttack when no bees are alive
func TestPlayerAttackEmptyHive(t *testing.T) {
	// Test attacking when no bees are alive
	game := NewGame()
	game.KillAllBees()

	// Capture stdout to verify the "no bees to attack" message
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	game.PlayerAttack()

	// Restore stdout and capture output
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Verify that the appropriate message is shown
	if !strings.Contains(output, "No bees left to attack!") {
		t.Errorf("Expected 'No bees left to attack!' message when attacking empty hive, got: %s", output)
	}

	// Verify game state remains consistent
	if len(game.GetAliveBees()) != 0 {
		t.Error("Should have no alive bees after attacking empty hive")
	}

	// Verify turn count doesn't increase when no attack occurs
	if game.Turns != 0 {
		t.Errorf("Expected turn count to remain 0 when no attack possible, got %d", game.Turns)
	}
}
