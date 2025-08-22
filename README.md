# BeesInTheTrap

A simple turn-based command-line game in Go.

## Project Structure

```
BeesInATrap/
├── cmd/beesinthetrap/     # Application entry point
│   └── main.go
├── internal/game/         # Game logic 
│   ├── bee.go
│   ├── player.go
│   ├── game.go
│   └── game_test.go
├── go.mod                 # Go module definition
└── README.md
```

## Setup

1. Build for Windows:
   ```
   go build -o beesinthetrap.exe ./cmd/beesinthetrap
   ```

## Run

```
go run ./cmd/beesinthetrap
```

## Test

```
go test ./internal/game
```
