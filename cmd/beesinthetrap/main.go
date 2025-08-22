package main

import (
	"fmt"

	"github.com/clearyalexandros/BeesInATrap/internal/game"
)

func main() {
	fmt.Println("Starting Bees in the Trap...")

	g := game.NewGame()
	g.Start()

	// Let's play!
	g.PlayGame()
}
