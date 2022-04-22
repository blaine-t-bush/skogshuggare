package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/gdamore/tcell"
)

func main() {
	// Seed randomizer.
	rand.Seed(time.Now().UTC().UnixNano())

	// Initialize tcell.
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
		player:   Actor{x: 5, y: 5},
		squirrel: Actor{x: 10, y: 10, destinationX: 20, destinationY: 20},
		border:   Border{0, w - 1, 0, h - 1, 1},
		trees:    map[int]*Tree{},
		exit:     false,
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

	// Randomly seed map with trees in various states.
	game.PopulateTrees(screen)

	// Perform first draw.
	game.Draw(screen)

	// Initialize game update ticker.
	ticker := time.NewTicker(TickRate * time.Millisecond)

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
			game.MoveActor(screen, ActorPlayer, 1, DirUp)
		case tcell.KeyRight:
			game.MoveActor(screen, ActorPlayer, 1, DirRight)
		case tcell.KeyDown:
			game.MoveActor(screen, ActorPlayer, 1, DirDown)
		case tcell.KeyLeft:
			game.MoveActor(screen, ActorPlayer, 1, DirLeft)
		case tcell.KeyRune:
			switch ev.Rune() {
			case rune(' '):
				game.Chop(screen, DirOmni)
			case rune('w'):
				game.Chop(screen, DirUp)
			case rune('d'):
				game.Chop(screen, DirRight)
			case rune('s'):
				game.Chop(screen, DirDown)
			case rune('a'):
				game.Chop(screen, DirLeft)
			}
		}
	case *tcell.EventResize:
		screen.Sync()
	}
	game.UpdateSquirrel(screen)
	//game.AddSeeds()
	game.GrowTrees()
}
