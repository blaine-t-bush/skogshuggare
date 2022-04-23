package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell"
)

func main() {
	// Initialize game state.
	game := Game{
		player:   Actor{position: Coordinate{x: 38, y: 5}, visionRadius: 100},
		squirrel: Actor{position: Coordinate{x: 10, y: 19}, visionRadius: 100},
		world:    readMap("kartor/flod.karta"),
		exit:     false,
	}

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

	// Randomly seed map with trees in various states.
	game.PopulateTrees(screen)
	game.PopulateGrass(screen)

	// Wait for Loop() goroutine to finish before moving on.
	var wg sync.WaitGroup
	wg.Add(1)
	go Ticker(&wg, screen, game)
	wg.Wait()
	screen.Fini()
}

func readMap(fileName string) World {
	filebuffer, err := ioutil.ReadFile(fileName)
	worldContent := make(map[Coordinate]interface{})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	filedata := string(filebuffer)
	data := bufio.NewScanner(strings.NewReader(filedata))
	data.Split(bufio.ScanLines)
	width := 0
	height := 0
	for data.Scan() {
		// Check if width needs to be updated. It's determined by the longest line.
		lineWidth := len(data.Text())
		if lineWidth > width {
			width = lineWidth
		}

		// Update the worldContent map according to special characters.
		for i := 0; i < lineWidth; i++ {
			switch data.Text()[i] {
			case '#': // wall
				worldContent[Coordinate{i, height}] = Object{'#', true, false, 255, 255, 255}
			case '~': // water
				worldContent[Coordinate{i, height}] = Object{'~', true, false, 10, 10, 200}
			case '=': // bridge
				worldContent[Coordinate{i, height}] = Object{'=', false, false, 150, 70, 0}
			}
		}

		// Increment the height once for each row.
		height++
	}

	return World{width, height, worldContent}
}

func Ticker(wg *sync.WaitGroup, screen tcell.Screen, game Game) {
	// Wait for this goroutine to finish before resuming main().
	defer wg.Done()

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
	// Listen for keyboard events for player actions,
	// or terminal resizing events to re-draw the screen.
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

	// Give the squirrel a destination if it doesn't alreasdy have one,
	// or update its destination if it's blocked.
	if (Coordinate{0, 0} == game.squirrel.destination) || game.IsBlocked(game.squirrel.destination) {
		game.squirrel.destination = game.GetRandomPlantableCoordinate()
	}

	// If squirrel is one move away from its destination, then it plants the seed at the destination,
	// i.e. one tile away, and then picks a new destination.
	// Otherwise, it just moves towards its current destination.
	if game.squirrel.IsAdjacentToDestination() {
		game.PlantSeed(game.squirrel.destination)
		game.squirrel.destination = game.GetRandomPlantableCoordinate()
	} else {
		var nextDirection int
		for {
			nextDirection = game.FindNextDirection(game.squirrel.position, game.squirrel.destination)
			if nextDirection == DirNone { // No path found, or on top of destination. Get a new one.
				game.squirrel.destination = game.GetRandomPlantableCoordinate()
			} else {
				break
			}
		}
		game.MoveActor(screen, ActorSquirrel, 1, nextDirection)
	}

	// Update trees.
	game.GrowTrees()
}
