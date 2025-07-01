package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"syscall"

	"golang.org/x/term"
)

type Bubble struct {
	X     int
	Y     int
	Text  string
	Alive bool
}

var width, height int

func clearScreen() {
	fmt.Print("\033[2J") // Clear screen
	fmt.Print("\033[H")  // Move cursor to top-left
}

func draw(bubbles []Bubble, fails int) {
	clearScreen()
	for _, b := range bubbles {
		if b.Alive {
			fmt.Printf("\033[%d;%dH%s", b.Y+1, b.X+1, b.Text)
		}
	}
	fmt.Printf("\033[%d;1HFails: %d", height, fails)
	fmt.Printf("\033[%d;1H", height+1) // move cursor below screen
}

func spawnBubble() Bubble {
	text := string(rune('a' + rand.Intn(26))) // random lowercase letter
	x := rand.Intn(width - 1)
	return Bubble{X: x, Y: 0, Text: text, Alive: true}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	width, height, _ = term.GetSize(int(os.Stdin.Fd()))
	height -= 2 // reserve lines for UI

	// Ctrl+C handler
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Println("\nExiting. Terminal restored.")
		os.Exit(0)
	}()

	var bubbles []Bubble
	var fails int
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	spawnTicker := time.NewTicker(1 * time.Second)
	defer spawnTicker.Stop()

	input := make(chan byte, 10)
	go func() {
		for {
			var b = make([]byte, 1)
			os.Stdin.Read(b)
			input <- b[0]
		}
	}()
	fmt.Println("Press Ctrl+C to exit. Type letters to pop bubbles.")

	for {
		select {
		case <-spawnTicker.C:
			bubbles = append(bubbles, spawnBubble())

		case <-ticker.C:
			// Move bubbles
			for i := range bubbles {
				if bubbles[i].Alive {
					bubbles[i].Y++
					if bubbles[i].Y >= height {
						bubbles[i].Alive = false
						fails++
					}
				}
			}

			// Remove dead bubbles (optional optimization)
			active := make([]Bubble, 0)
			for _, b := range bubbles {
				if b.Alive {
					active = append(active, b)
				}
			}
			bubbles = active

			draw(bubbles, fails)

			if fails >= 5 {
				fmt.Println("Game Over! You failed 5 times.")
				return
			}

		case key := <-input:
			for i := range bubbles {
				if bubbles[i].Alive && bubbles[i].Text[0] == key {
					bubbles[i].Alive = false
					break
				}
			}
		}
	}
}
