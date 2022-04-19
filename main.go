package main

import (
	"fmt"
	"os"
	"sync"

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
		trees: []Tree{
			{10, 10},
			{15, 5},
			{16, 8},
			{20, 15},
			{25, 10},
		},
	}

	// Wait for Loop() goroutine to finish before moving on.
	var wg sync.WaitGroup
	wg.Add(1)
	go Loop(&wg, screen, game)
	wg.Wait()
	screen.Fini()
}

func Loop(wg *sync.WaitGroup, screen tcell.Screen, game Game) {
	defer wg.Done()
	// Perform first draw.
	game.DrawPlayer(screen)
	game.border.Draw(screen)
	for _, tree := range game.trees {
		tree.Draw(screen)
	}

	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				return
			case tcell.KeyUp:
				game.MovePlayer(screen, 1, 0)
			case tcell.KeyRight:
				game.MovePlayer(screen, 1, 1)
			case tcell.KeyDown:
				game.MovePlayer(screen, 1, 2)
			case tcell.KeyLeft:
				game.MovePlayer(screen, 1, 3)
			}
		case *tcell.EventResize:
			screen.Sync()
		}
	}
}
