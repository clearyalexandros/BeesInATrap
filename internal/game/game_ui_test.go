package game

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

// Test Start function welcome message and initial display
func TestStartGameWelcomeMessage(t *testing.T) {
	game := NewGame()

	// Capture stdout to verify output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Test that Start executes without panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Start() panicked: %v", r)
		}
	}()

	game.Start()

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify expected content in output
	expectedPhrases := []string{
		"Welcome to Bees in the Trap!",
		"Your mission: Destroy the hive",
		"Type 'hit' to attack",
		"Player HP: 100/100",
		"Queens: 1",
		"Workers: 5",
		"Drones: 25",
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(output, phrase) {
			t.Errorf("Expected Start() output to contain '%s', but it didn't. Output: %s", phrase, output)
		}
	}
}

// Test PlayGame with mocked input - Hit command
func TestPlayGameHitCommand(t *testing.T) {
	game := NewGame()

	// Mock stdin with "hit" followed by "quit"
	input := "hit\nquit\n"
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r

	go func() {
		defer w.Close()
		w.Write([]byte(input))
	}()

	// Capture stdout to avoid clutter
	oldStdout := os.Stdout
	os.Stdout, _, _ = os.Pipe()

	// Test that PlayGame executes without panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PlayGame() panicked: %v", r)
		}
		// Restore stdin/stdout
		os.Stdin = oldStdin
		os.Stdout = oldStdout
	}()

	game.PlayGame()

	// Verify that a turn was taken
	if game.Turns == 0 {
		t.Error("Expected at least one turn to be taken after 'hit' command")
	}
}

// Test PlayGame with auto mode
func TestPlayGameAutoMode(t *testing.T) {
	game := NewGame()

	// Mock stdin with "auto" command
	input := "auto\n"
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r

	go func() {
		defer w.Close()
		w.Write([]byte(input))
	}()

	// Capture stdout to avoid clutter
	oldStdout := os.Stdout
	os.Stdout, _, _ = os.Pipe()

	// Set a timeout to prevent infinite auto-play
	done := make(chan bool, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("PlayGame() in auto mode panicked: %v", r)
			}
			done <- true
		}()

		game.PlayGame()
	}()

	// Wait for either completion or timeout
	select {
	case <-done:
		// Game completed normally
	case <-time.After(2 * time.Second):
		// Timeout - force game end by killing all bees
		game.KillAllBees()
		<-done // Wait for goroutine to finish
	}

	// Restore stdin/stdout
	os.Stdin = oldStdin
	os.Stdout = oldStdout

	// Verify auto mode was activated
	if !game.AutoMode {
		t.Error("Expected AutoMode to be true after 'auto' command")
	}

	// Verify turns were taken
	if game.Turns == 0 {
		t.Error("Expected turns to be taken in auto mode")
	}
}

// Test PlayGame with invalid commands
func TestPlayGameInvalidCommands(t *testing.T) {
	game := NewGame()

	// Mock stdin with invalid commands followed by quit
	input := "invalid\nwrong\nbad\nquit\n"
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r

	go func() {
		defer w.Close()
		w.Write([]byte(input))
	}()

	// Capture stdout to check for error messages
	oldStdout := os.Stdout
	captureR, captureW, _ := os.Pipe()
	os.Stdout = captureW

	// Test that PlayGame handles invalid commands gracefully
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PlayGame() panicked on invalid commands: %v", r)
		}
	}()

	game.PlayGame()

	// Restore stdin/stdout and read captured output
	captureW.Close()
	os.Stdin = oldStdin
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, captureR)
	output := buf.String()

	// Verify error message for invalid commands
	if !strings.Contains(output, "Invalid command") {
		t.Error("Expected error message for invalid commands")
	}

	// Game should not have progressed with invalid commands
	if game.Turns != 0 {
		t.Error("Expected no turns to be taken with only invalid commands")
	}
}

// Test PlayGame quit command
func TestPlayGameQuitCommand(t *testing.T) {
	game := NewGame()

	// Mock stdin with immediate quit
	input := "quit\n"
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r

	go func() {
		defer w.Close()
		w.Write([]byte(input))
	}()

	// Capture stdout to check for quit message
	oldStdout := os.Stdout
	captureR, captureW, _ := os.Pipe()
	os.Stdout = captureW

	// Test that PlayGame handles quit gracefully
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PlayGame() panicked on quit: %v", r)
		}
	}()

	game.PlayGame()

	// Restore stdin/stdout and read captured output
	captureW.Close()
	os.Stdin = oldStdin
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, captureR)
	output := buf.String()

	// Verify quit message
	if !strings.Contains(output, "Thanks for playing!") {
		t.Error("Expected quit message to be displayed")
	}

	// Game should not have progressed
	if game.Turns != 0 {
		t.Error("Expected no turns to be taken when quitting immediately")
	}
}

// Test PlayGame with EOF (scanner fails)
func TestPlayGameEOF(t *testing.T) {
	game := NewGame()

	// Mock stdin that closes immediately (EOF)
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	w.Close() // Close immediately to simulate EOF

	// Capture stdout to avoid clutter
	oldStdout := os.Stdout
	os.Stdout, _, _ = os.Pipe()

	// Test that PlayGame handles EOF gracefully
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PlayGame() panicked on EOF: %v", r)
		}
		// Restore stdin/stdout
		os.Stdin = oldStdin
		os.Stdout = oldStdout
	}()

	game.PlayGame()

	// Should exit gracefully without panic
}

// Test PlayGame complete game flow
func TestPlayGameCompleteFlow(t *testing.T) {
	game := NewGame()

	// Kill ALL bees to make game end immediately
	game.KillAllBees()

	// Mock stdin with hit command
	input := "hit\n"
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r

	go func() {
		defer w.Close()
		w.Write([]byte(input))
	}()

	// Capture stdout to avoid clutter
	oldStdout := os.Stdout
	os.Stdout, _, _ = os.Pipe()

	// Test complete game flow
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PlayGame() complete flow panicked: %v", r)
		}
		// Restore stdin/stdout
		os.Stdin = oldStdin
		os.Stdout = oldStdout
	}()

	// Game should end immediately since all bees are dead
	game.PlayGame()

	// Game should have ended
	if !game.IsGameOver() {
		t.Error("Expected game to be over after complete flow")
	}
}
