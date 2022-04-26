package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell"
)

func main() {
	// Attempt to get map name from command line args.
	var mapName string
	var visionRadius int
	if len(os.Args) <= 1 { // Make sure there are arguments before accesing slices
		mapName = "skog"
		visionRadius = 20
	} else {
		if len(os.Args[1:]) >= 1 {
			mapName = os.Args[1]
		} else {
			// Couldn't parse map name from command line. Using default map.
			mapName = "skog"
		}

		// Attempt to get vision radius from command line args.

		if len(os.Args[2:]) >= 1 {
			visionRadius, _ = strconv.Atoi(os.Args[2])
		} else {
			// Couldn't parse map name from command line. Using default map.
			visionRadius = 20
		}

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

	// Draw and handle menu inputs before initializing and drawing the game itself
	titleMenu := GenerateTitleMenu() // Generate title menu
	var twg sync.WaitGroup
	twg.Add(1)
	go TitleMenuHandler(&twg, screen, &titleMenu)
	twg.Wait()

	// Read map to initialize game state.
	worldContent, playerPosition, squirrelPosition := ReadMap("kartor/" + mapName + ".karta")
	game := Game{
		player:   Actor{position: playerPosition, visionRadius: visionRadius, score: 0},
		squirrel: Actor{position: squirrelPosition, visionRadius: 100, score: 0},
		world:    worldContent,
		menu:     Menu{15, 5, Coordinate{0, 0}, []string{}},
		exit:     false,
	}

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

func TitleMenuHandler(wg *sync.WaitGroup, screen tcell.Screen, titleMenu *TitleMenu) { // TODO make sure variables are not changed at the same time w/ mutex or channels
	defer wg.Done()

	// Initialize game menu update ticker.
	ticker := time.NewTicker(TickRate * time.Millisecond)

	// Start the input handler
	go titleMenu.InputHandler(screen)

	for range ticker.C {
		titleMenu.Update(screen)
		titleMenu.Draw(screen)
		if titleMenu.exit {
			screen.Clear()
			return
		}
	}
}

func ReadMap(fileName string) (World, Coordinate, Coordinate) {
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
	var playerPosition, squirrelPosition Coordinate
	for data.Scan() {
		// Check if width needs to be updated. It's determined by the longest line.
		lineWidth := len(data.Text())
		if lineWidth > width {
			width = lineWidth
		}

		// Update the worldContent map according to special characters.
		for i := 0; i < lineWidth; i++ {
			switch data.Text()[i] {
			case MapPlayer:
				playerPosition = Coordinate{i, height}
			case MapSquirrel:
				squirrelPosition = Coordinate{i, height}
			case MapWall:
				worldContent[Coordinate{i, height}] = Object{KeyWall, true, false}
			case MapWaterLight:
				worldContent[Coordinate{i, height}] = Object{KeyWaterLight, true, false}
			case MapWaterHeavy:
				worldContent[Coordinate{i, height}] = Object{KeyWaterHeavy, true, false}
			}
		}

		// Increment the height once for each row.
		height++
	}

	_borders := make(map[Coordinate]int)

	for c := range worldContent {
		if border, isBorder := IsBorder(width, height, c); isBorder {
			_borders[c] = border
		}
	}

	return World{width, height, _borders, worldContent}, playerPosition, squirrelPosition
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
				game.Chop(screen, DirOmni, 1)
			case rune('w'):
				game.Chop(screen, DirUp, 1)
			case rune('d'):
				game.Chop(screen, DirRight, 1)
			case rune('s'):
				game.Chop(screen, DirDown, 1)
			case rune('a'):
				game.Chop(screen, DirLeft, 1)
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
