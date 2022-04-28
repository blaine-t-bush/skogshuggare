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
	// Attempt to get vision radius from command line args.
	visionRadius := 100
	if len(os.Args) >= 2 { // Make sure there are arguments before accessing slices
		visionRadius, _ = strconv.Atoi(os.Args[1])
	}

	// Seed randomizer.
	rand.Seed(time.Now().UTC().UnixNano())

	// Initialize game state.
	var game Game
	var err error

	// Initialize tcell.
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	game.screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err = game.screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	// Set default style and clear terminal.
	game.screen.SetStyle(tcell.StyleDefault)
	game.screen.Clear()

	// Draw and handle menu inputs before initializing and drawing the game itself
	titleMenu := GenerateTitleMenu() // Generate title menu
	var twg sync.WaitGroup
	twg.Add(1)
	go TitleMenuHandler(&twg, game.screen, &titleMenu)
	twg.Wait()

	// Read map to initialize game state.
	mapName := titleMenu.selectedMap
	worldContent, playerPosition, squirrelPositions := ReadMap("kartor/" + mapName)
	squirrels := make(map[int]*Actor)
	for index, position := range squirrelPositions {
		squirrels[index] = &Actor{position: position, visionRadius: 100, score: 0, hitPointsCurrent: MaxHitPointsSquirrel, hitPointsMax: MaxHitPointsSquirrel}
	}
	game.player = Actor{position: playerPosition, visionRadius: visionRadius, score: 0, hitPointsCurrent: MaxHitPointsPlayer, hitPointsMax: MaxHitPointsPlayer}
	game.squirrels = squirrels
	game.world = worldContent
	game.menu = Menu{15, 5, Coordinate{0, 0}, []string{}}
	game.exit = false

	// Randomly seed map with trees in various states.
	game.PopulateTrees()
	game.PopulateGrass()

	// Create a mutex so Ticker and AnimationHandler can both use the game map
	// "concurrently" without colliding. To be more specific, the mutex prevents
	// the two goroutines from doing exactly concurrent operations.
	mutex := &sync.Mutex{}
	// Wait for Loop() goroutine to finish before moving on.
	var wg sync.WaitGroup
	wg.Add(2)
	go game.StateTicker(&wg, mutex)
	go game.AnimationTicker(&wg, mutex)
	wg.Wait()
	game.screen.Fini()

	fmt.Println("Game over. Final score:", game.player.score)
	fmt.Println("")
}

func TitleMenuHandler(wg *sync.WaitGroup, screen tcell.Screen, titleMenu *TitleMenu) { // TODO make sure variables are not changed at the same time w/ mutex or channels
	defer wg.Done()

	// Initialize game menu update ticker.
	ticker := time.NewTicker(StateTickDuration * time.Millisecond)

	// Start the input handler
	go titleMenu.InputHandler(screen)
	// Start title menu animation handler
	go titleMenu.AnimationHandler()

	for range ticker.C {
		titleMenu.Update(screen)
		titleMenu.Draw(screen)
		if titleMenu.exit {
			screen.Clear()
			return
		}
	}
}

func ReadMap(fileName string) (World, Coordinate, []Coordinate) {
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
	var playerPosition Coordinate
	var squirrelPositions []Coordinate
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
				squirrelPositions = append(squirrelPositions, Coordinate{i, height})
			case MapWall:
				worldContent[Coordinate{i, height}] = &StaticObject{KeyWall, true, false, false}
			case MapWater:
				worldContent[Coordinate{i, height}] = &AnimatedObject{KeyWater, GetRandomAnimationStage(KeyWater), true, false, false}
			case MapFire:
				worldContent[Coordinate{i, height}] = &Fire{GetRandomAnimationStage(KeyFire), Coordinate{i, height}, 0}
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

	return World{width, height, _borders, worldContent}, playerPosition, squirrelPositions
}

func (game *Game) StateTicker(wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()

	// Initialize game update ticker.
	ticker := time.NewTicker(StateTickDuration * time.Millisecond)

	// Update game state and re-draw on every tick.
	for range ticker.C {
		mutex.Lock()
		game.Draw()
		mutex.Unlock()
		game.StateUpdate(mutex)
		if game.exit {
			wg.Done()
			return
		}
	}
}

func (game *Game) StateUpdate(mutex *sync.Mutex) {
	// Listen for keyboard events for player actions,
	// or terminal resizing events to re-draw the screen.
	ev := game.screen.PollEvent()
	mutex.Lock()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape:
			game.exit = true
			return
		case tcell.KeyUp:
			game.MovePlayer(1, DirUp)
		case tcell.KeyRight:
			game.MovePlayer(1, DirRight)
		case tcell.KeyDown:
			game.MovePlayer(1, DirDown)
		case tcell.KeyLeft:
			game.MovePlayer(1, DirLeft)
		case tcell.KeyRune:
			switch ev.Rune() {
			case rune('q'):
				game.Chop(DirOmni, 1)
			case rune('w'):
				game.Chop(DirUp, 1)
			case rune('d'):
				game.Chop(DirRight, 1)
			case rune('s'):
				game.Chop(DirDown, 1)
			case rune('a'):
				game.Chop(DirLeft, 1)
			case rune('Q'):
				game.Dig(DirOmni)
			case rune('W'):
				game.Dig(DirUp)
			case rune('D'):
				game.Dig(DirRight)
			case rune('S'):
				game.Dig(DirDown)
			case rune('A'):
				game.Dig(DirLeft)
			}
		}
	case *tcell.EventResize:
		game.screen.Sync()
	}

	// Give the squirrel a destination if it doesn't alreasdy have one,
	// or update its destination if it's blocked.
	// FIXME determine why squirrels sometimes stop even when there seem to be nearby available plantable coordinates
	// FIXME move to function
	for key, squirrel := range game.squirrels {
		if (Coordinate{0, 0} == squirrel.destination) || game.IsPathBlocked(squirrel.destination) {
			squirrel.destination = game.GetRandomPlantableCoordinate()
			squirrel.path = game.FindPath(squirrel.position, squirrel.destination)
		}

		// If squirrel is one move away from its destination, then it plants the seed at the destination,
		// i.e. one tile away, and then picks a new destination.
		// Otherwise, it just moves towards its current destination.
		if squirrel.IsAdjacentToDestination() && !game.IsPathBlocked(squirrel.destination) {
			game.PlantSeed(squirrel.destination)
			squirrel.destination = game.GetRandomPlantableCoordinate()
		} else {
			nextDirection := game.FindNextDirection(key)
			if nextDirection == DirNone { // No path found, or on top of destination. Get a new one.
				squirrel.destination = game.GetRandomPlantableCoordinate()
			}
			game.MoveSquirrel(1, nextDirection, key)
			game.UpdatePath(key)
		}
	}

	// Update trees.
	game.GrowTrees()

	// Update fire.
	game.UpdateFire()
	game.CheckFireDamage()
	mutex.Unlock()
}
