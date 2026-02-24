# ⌨️ Terminal Typing Game

A fast-paced, colorful terminal typing game where bubbles fall from the top of the screen. You must type their respective characters to pop them before they hit the bottom! The longer you play, the harder it gets.

---

## 🎮 Gameplay Features

- **Falling Bubbles**: Bubbles containing random characters (a-z, A-Z, 0-9) fall from the top of your terminal.
- **Score & Streak System**: Popping a bubble earns you points. Consecutive pops build a streak that gives you bonus points! A single miss resets your streak.
- **Dynamic Difficulty (Levels)**: Every 150 points, you level up. The game speeds up, with bubbles spawning faster and falling quicker.
- **Lives System**: You start with 5 lives (❤️). Letting a bubble reach the bottom costs a life and causes the screen to flash red. When you hit 0 lives, it's Game Over.
- **Colorful UI**: Bubbles spawn in vibrant ANSI colors. A status bar at the bottom cleanly tracks your Progress, Score, Lives, Level, and Streak.
- **Clean Exit**: Press `Esc` or `Ctrl+C` at any time to safely exit and restore your terminal cursor.

---

## 🚀 Getting Started

### ✅ Prerequisites

- [Go](https://golang.org/dl/) 1.17+
- A terminal that supports ANSI escape codes (most modern terminals do).

### 🛠 Setup & Run

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/terminal_typing_game.git
   cd terminal_typing_game
   ```

2. To run the game directly:
   ```bash
   go run main.go
   ```

3. To build the executable:
   ```bash
   go build -o terminal_typing_game
   ./terminal_typing_game
   ```

Enjoy the game, and try to beat your high score!
