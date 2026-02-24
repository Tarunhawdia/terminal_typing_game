package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/term"
)

type Bubble struct {
	X     int
	Y     int
	Text  string
	Color string
	Alive bool
}

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

var colors = []string{ColorRed, ColorGreen, ColorYellow, ColorBlue, ColorPurple, ColorCyan, ColorWhite}

var width, height int

func clearScreen() {
	fmt.Print("\033[2J") // Clear screen
	fmt.Print("\033[H")  // Move cursor to top-left
}

func draw(bubbles []Bubble, lives int, score int, level int, streak int) {
	clearScreen()
	for _, b := range bubbles {
		if b.Alive {
			fmt.Printf("\033[%d;%dH%s%s%s", b.Y+1, b.X+1, b.Color, b.Text, ColorReset)
		}
	}

	// Draw UI bar at the bottom
	uiY := height

	// Format Lives
	livesStr := ""
	for i := 0; i < lives; i++ {
		livesStr += "❤️ "
	}

	header := fmt.Sprintf("\033[%d;1H%s| Score: %d | Level: %d | Streak: %d | Lives: %s%s", uiY, ColorCyan, score, level, streak, livesStr, ColorReset)
	fmt.Print(header)
	fmt.Printf("\033[%d;1H", uiY+1) // move cursor below screen
}

func spawnBubble() Bubble {
	charsets := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	text := string(charsets[rand.Intn(len(charsets))])

	x := 1
	if width > 2 {
		x = rand.Intn(width - 2)
		if x <= 0 {
			x = 1
		}
	}

	color := colors[rand.Intn(len(colors))]
	return Bubble{X: x, Y: 0, Text: text, Color: color, Alive: true}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	fmt.Print("\033[?25l")       // Hide cursor
	defer fmt.Print("\033[?25h") // Show cursor on exit

	width, height, err = term.GetSize(int(os.Stdin.Fd()))
	if err != nil || width == 0 || height == 0 {
		width = 80
		height = 24
	}
	height -= 2 // reserve lines for UI

	// Ctrl+C handler
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Print("\033[?25h\033[2J\033[H") // Show cursor and clear screen
		fmt.Println("Exiting. Terminal restored.")
		os.Exit(0)
	}()

	var bubbles []Bubble
	lives := 5
	score := 0
	level := 1
	streak := 0

	moveSpeed := 150 * time.Millisecond
	spawnSpeed := 1200 * time.Millisecond

	ticker := time.NewTicker(moveSpeed)
	defer ticker.Stop()

	spawnTicker := time.NewTicker(spawnSpeed)
	defer spawnTicker.Stop()

	input := make(chan byte, 10)
	go func() {
		for {
			var b = make([]byte, 1)
			_, err := os.Stdin.Read(b)
			if err == nil {
				input <- b[0]
			}
		}
	}()

	clearScreen()
	titleMsg := " TERMINAL TYPING GAME "
	fmt.Printf("\033[%d;%dH%s%s%s", height/2-2, width/2-len(titleMsg)/2, ColorGreen, titleMsg, ColorReset)
	startMsg := "Press any key to start typing..."
	fmt.Printf("\033[%d;%dH%s", height/2, width/2-len(startMsg)/2, startMsg)
	<-input

	for {
		select {
		case <-spawnTicker.C:
			bubbles = append(bubbles, spawnBubble())

		case <-ticker.C:
			// Move bubbles
			for i := range bubbles {
				if bubbles[i].Alive {
					bubbles[i].Y++
					if bubbles[i].Y >= height-1 {
						bubbles[i].Alive = false
						lives--
						streak = 0 // Reset streak on miss
						if lives > 0 {
							// Flash screen red briefly
							fmt.Print("\033[41m") // Red background
							draw(bubbles, lives, score, level, streak)
							time.Sleep(50 * time.Millisecond)
							fmt.Print("\033[0m") // Reset
							draw(bubbles, lives, score, level, streak)
						}
					}
				}
			}

			// Remove dead bubbles
			active := make([]Bubble, 0)
			for _, b := range bubbles {
				if b.Alive {
					active = append(active, b)
				}
			}
			bubbles = active

			draw(bubbles, lives, score, level, streak)

			if lives <= 0 {
				term.Restore(int(os.Stdin.Fd()), oldState)
				fmt.Print("\033[?25h\033[2J\033[H") // Show cursor and clear screen

				gameOverMsg := "GAME OVER!"
				fmt.Printf("\033[%d;%dH%s%s%s\n", height/2-2, width/2-len(gameOverMsg)/2, ColorRed, gameOverMsg, ColorReset)

				scoreMsg := fmt.Sprintf("Final Score: %d", score)
				fmt.Printf("\033[%d;%dH%s%s%s\n", height/2, width/2-len(scoreMsg)/2, ColorGreen, scoreMsg, ColorReset)

				levelMsg := fmt.Sprintf("Level Reached: %d", level)
				fmt.Printf("\033[%d;%dH%s%s%s\n", height/2+1, width/2-len(levelMsg)/2, ColorYellow, levelMsg, ColorReset)

				fmt.Printf("\033[%d;0H\n", height) // Move to bottom
				return
			}

		case key := <-input:
			// Handle inputs like Ctrl+C manually if it passes through raw mode as \x03
			// \x1B is Esc
			if key == 3 || key == 27 {
				term.Restore(int(os.Stdin.Fd()), oldState)
				fmt.Print("\033[?25h\033[2J\033[H")
				fmt.Println("Exiting...")
				return
			}

			popped := false
			for i := range bubbles {
				if bubbles[i].Alive && bubbles[i].Text[0] == key {
					bubbles[i].Alive = false
					popped = true
					streak++
					score += 10 + streak // Bonus for streak
					break                // Only pop one bubble at a time (lowest first)
				}
			}

			if popped {
				// Level up every 150 points
				newLevel := (score / 150) + 1
				if newLevel > level {
					level = newLevel
					// Increase difficulty
					if moveSpeed > 40*time.Millisecond {
						moveSpeed -= 15 * time.Millisecond
						ticker.Reset(moveSpeed)
					}
					if spawnSpeed > 300*time.Millisecond {
						spawnSpeed -= 150 * time.Millisecond
						spawnTicker.Reset(spawnSpeed)
					}
				}
				draw(bubbles, lives, score, level, streak)
			} else {
				// Missed a keystroke
				streak = 0
				draw(bubbles, lives, score, level, streak)
			}
		}
	}
}
