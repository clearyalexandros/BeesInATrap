# BeesInTheTrap

A turn-based command-line game in Go where you battle against a hive of bees! Destroy all the bees before they sting you to death.

## 🎮 How to Play

### Game Objective

Your mission is simple: **Destroy the entire hive before the bees sting you to death!**

### Starting the Game

When you run the game, you'll see:

- Your health: **100/100 HP**
- The hive composition: **1 Queen, 5 Workers, 25 Drones** (31 bees total)
- A command prompt asking for your action

### User Commands

The game accepts three simple commands:

| Command | Description |
|---------|-------------|
| `hit` | Attack the hive - you'll target a random bee |
| `auto` | Switch to automatic mode - the game plays itself |
| `quit` | Exit the game immediately |

### Game Flow

#### 1. **Player Turn**

- Type `hit` to attack the hive
- You have a **15% chance to miss** completely
- If you hit, you'll damage a random bee:
  - **Queen**: Takes 10 damage (100 HP total)
  - **Worker**: Takes 25 damage (75 HP total)
  - **Drone**: Takes 30 damage (60 HP total)
- **Special Rule**: Killing the Queen instantly eliminates all remaining bees!

#### 2. **Bees Turn**

- All living bees "think" simultaneously using goroutines
- Each bee has a **20% chance to miss** their attack
- If they hit, damage varies by bee type:
  - **Queen**: 10 damage per sting 🩸
  - **Worker**: 5 damage per sting ⚡
  - **Drone**: 1 damage per sting 🔸
- Real-time damage alerts show your health status

#### 3. **Victory Conditions**

- **You Win**: Eliminate all bees (or kill the Queen)
- **You Lose**: Your HP reaches 0

### Example Gameplay Session

```text
Welcome to Bees in the Trap!
Your mission: Destroy the hive before the bees sting you to death!

=== Game Status ===
Player HP: 100/100
Alive Bees:
  Queens: 1
  Workers: 5  
  Drones: 25
Turns: 0

Enter command (hit/auto/quit): hit

--- Turn 1: Player Turn ---
Direct Hit! You attacked a Drone bee!
The Drone bee took 30 damage and has 30 HP remaining.

--- Turn 1: Bees Turn ---
🧠 Bees consulted for 1.2s total...
Sting! You just got stung by a Worker bee!
You took 5 damage and now have 95 HP remaining.
⚡ Damage Alert: -5 HP | Turn 1 | Player: 95/100 (95.0%) | Bees: 31

Enter command (hit/auto/quit): auto
Switching to auto mode...
```

### Auto Mode

- Type `auto` to let the computer play automatically
- The game continues until victory or defeat
- Perfect for demonstrations or when you want to watch the AI battle!

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

### Cross-Platform Compilation

Build executable binaries for both Windows and Linux:

1. **Build for Windows:**

   ```bash
   # Unix/Linux/macOS
   GOOS=windows GOARCH=amd64 go build -o beesinthetrap.exe ./cmd/beesinthetrap
   ```

   ```powershell
   # Windows PowerShell
   $env:GOOS="windows"; $env:GOARCH="amd64"; go build -o beesinthetrap.exe ./cmd/beesinthetrap
   ```

2. **Build for Linux:**

   ```bash
   # Unix/Linux/macOS
   GOOS=linux GOARCH=amd64 go build -o beesinthetrap ./cmd/beesinthetrap
   ```

   ```powershell
   # Windows PowerShell
   $env:GOOS="linux"; $env:GOARCH="amd64"; go build -o beesinthetrap ./cmd/beesinthetrap
   ```

## Run

```bash
# Run with default settings
go run ./cmd/beesinthetrap

# Run with custom configuration
go run ./cmd/beesinthetrap --player-hp 150 --player-miss 0.10 --bees-miss 0.30

# Create a custom hive composition
go run ./cmd/beesinthetrap --queens 2 --workers 10 --drones 50

# Easy mode (high player HP, low miss chance, slow bees)
go run ./cmd/beesinthetrap --player-hp 200 --player-miss 0.05 --bees-miss 0.40

# Hard mode (low player HP, high miss chance, fast auto mode)
go run ./cmd/beesinthetrap --player-hp 50 --player-miss 0.25 --bees-miss 0.10 --auto-delay 200

# See all configuration options
go run ./cmd/beesinthetrap --help
```

### Configuration Flags

| Flag | Description | Default | Range |
|------|-------------|---------|-------|
| `--player-hp` | Starting health points for the player | 100 | > 0 |
| `--player-miss` | Player miss chance | 0.15 (15%) | 0.0-1.0 |
| `--bees-miss` | Bees miss chance | 0.20 (20%) | 0.0-1.0 |
| `--auto-delay` | Auto mode delay in milliseconds | 500 | ≥ 0 |
| `--queens` | Number of Queen bees in the hive | 1 | ≥ 0 |
| `--workers` | Number of Worker bees in the hive | 5 | ≥ 0 |
| `--drones` | Number of Drone bees in the hive | 25 | ≥ 0 |
| `--help` | Show help information | - | - |

## Test

```bash
go test ./internal/game
```

## Docker Deployment

### Build and Run

```bash
# Build the image
docker build -t beesinthetrap:latest .

# Run interactively
docker run -it --rm beesinthetrap:latest
```

### Docker Compose

```bash
# Run with interactive terminal (recommended)
docker-compose run --rm beesinthetrap

# Auto-play mode
echo "auto" | docker-compose run --rm beesinthetrap
```

### Features

- ✅ **18.5MB image** - Optimized multi-stage build
- ✅ **Secure** - Non-root user execution  
- ✅ **Interactive** - Full terminal support for gameplay
