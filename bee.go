package main

type BeeType int

const (
	Queen BeeType = iota
	Worker
	Drone
)

type Bee struct {
	Type  BeeType
	HP    int
	Alive bool
}

// TODO: Add bee actions and logic
