# BeesInTheTrap

A turn-based command-line game in Go that demonstrates concurrent programming with goroutines and channels.

## Concurrency Features

This game showcases Go's concurrency.

### 🐝 **Concurrent Bee Decision Making**

Each bee in the hive makes attack decisions simultaneously using individual goroutines. When it's the bees' turn, all 31 bees "think" in parallel about whether to attack or miss, with realistic thinking times based on bee type:

- Queen bees: 50-150ms (strategic decisions)
- Worker bees: 20-80ms (moderate thinking)
- Drone bees: 10-50ms (quick reactions)

### 📊 **Real-Time Damage Monitoring**

A background goroutine listens for damage events through a buffered channel. Every time the player takes damage, it instantly displays live statistics including health percentage, bee count, and turn information with visual indicators:

- 🔸 Light damage (1-4 HP)
- ⚡ Medium damage (5-9 HP)
- 🩸 Heavy damage (10+ HP)

### 🔄 **Event-Driven Architecture**

The game uses channels for non-blocking communication between goroutines, ensuring smooth gameplay while background processes handle monitoring, statistics, and concurrent bee behavior without interrupting the main game flow.

## Project Structure

```text
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

   ```bash
   go build -o beesinthetrap.exe ./cmd/beesinthetrap
   ```

## Run

```bash
go run ./cmd/beesinthetrap
```

## Test

```bash
go test ./internal/game
```
