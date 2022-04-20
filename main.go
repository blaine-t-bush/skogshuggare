package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gdamore/tcell"
)

func main() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err = screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	// Set default style and clear terminal.
	screen.SetStyle(tcell.StyleDefault)
	screen.Clear()

	// Initialize game state.
	w, h := screen.Size()
	game := Game{
		player: Player{5, 5},
		border: Border{0, w - 1, 0, h - 1, 1},
		trees: map[int]*Tree{
			0: {6, 5, 2},
			1: {3, 5, 2},
			2: {10, 10, 2},
			3: {15, 5, 2},
			4: {16, 8, 2},
			5: {20, 15, 2},
			6: {25, 10, 2},
		},
		exit: false,
	}

	// Wait for Loop() goroutine to finish before moving on.
	var wg sync.WaitGroup
	wg.Add(1)
	go Ticker(&wg, screen, game)
	wg.Wait()
	screen.Fini()
}

func Ticker(wg *sync.WaitGroup, screen tcell.Screen, game Game) {
	// Wait for this goroutine to finish before resuming main().
	defer wg.Done()

	// Perform first draw.
	game.Draw(screen)

	// Initialize game update ticker.
	ticker := time.NewTicker(30 * time.Millisecond)

	// Update game state and re-draw on every tick.
	for range ticker.C {
		game.Update(screen)
		game.Draw(screen)
		if game.exit {
			return
		}
	}
}

func (game *Game) Update(screen tcell.Screen) {
	ev := screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape:
			game.exit = true
			return
		case tcell.KeyUp:
			game.MovePlayer(screen, 1, 0)
		case tcell.KeyRight:
			game.MovePlayer(screen, 1, 1)
		case tcell.KeyDown:
			game.MovePlayer(screen, 1, 2)
		case tcell.KeyLeft:
			game.MovePlayer(screen, 1, 3)
		case tcell.KeyRune:
			switch ev.Rune() {
			case rune(' '):
				game.Chop(screen, -1)
			case rune('w'):
				game.Chop(screen, 0)
			case rune('d'):
				game.Chop(screen, 1)
			case rune('s'):
				game.Chop(screen, 2)
			case rune('a'):
				game.Chop(screen, 3)
			}
		}
	case *tcell.EventResize:
		screen.Sync()
	}
}
