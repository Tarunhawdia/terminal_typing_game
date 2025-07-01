# ⌨️ Typing Bubble Game (Terminal Edition)

A fast-paced terminal typing game where bubbles fall from the top of the screen and you must type their letters to pop them before they hit the bottom. Miss 5 and it's game over.

---

## 🎮 Gameplay

- Bubbles with random letters fall from the top.
- Type the correct letter to pop a bubble.
- If a bubble reaches the bottom, it's counted as a fail.
- After 5 fails → **Game Over**.
- Press `q` at any time to quit.

---

## 🚀 Getting Started

### ✅ Prerequisites

- [Go](https://golang.org/dl/) 1.17+
- Terminal that supports ANSI escape codes (most do)

### 🛠 Setup

```bash
git clone https://github.com/yourusername/typing-bubble-game.git
cd typing-bubble-game
go mod init typing-bubble-game
go get golang.org/x/term
go run main.go
