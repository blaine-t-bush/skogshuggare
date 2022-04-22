package main

import (
	"math/rand"

	"github.com/gdamore/tcell"
)

type Player struct {
	x             int // Player x-coordinate
	y             int // Player y-coordinate
	vision_radius int // Player vision radius
}

type Border struct {
	x1 int // Left border boundary x-coordinate
	x2 int // Right border boundary x-coordinate
	y1 int // Top border boundary y-coordinate
	y2 int // Bottom border boundary y-coordinate
	t  int // Border thickness in characters
}

type Tree struct {
	x     int // Trunk left corner x-coordinate
	y     int // Trunk left corner y-coordinate
	state int // See constants
}

type Object struct {
	char       rune
	collidable bool
}

const (
	// Game parameters
	TickRate = 30 // Milliseconds between ticks
	// Directions
	DirUp    = 0
	DirRight = 1
	DirDown  = 2
	DirLeft  = 3
	DirOmni  = 4
	// Living tree states
	TreeStateSeed    = 0
	TreeStateSapling = 1
	TreeStateAdult   = 2
	// Harvested tree states
	TreeStateRemoved   = 10
	TreeStateStump     = 11
	TreeStateTrunk     = 12
	TreeStateStumpling = 13
	// Growth chances (per game tick)
	GrowthChanceSeed    = 0.010 // Seed to sapling
	GrowthChanceSapling = 0.005 // Sapling to adult
	SeedCreationChance  = 0.005 // Seed spawning
	SeedCreationMax     = 3     // Maximum number of seeds to create per tick
)

type Coordinate struct {
	x int
	y int
}

type World struct {
	width   int
	height  int
	content map[Coordinate]interface{}
}

type Game struct {
	player Player
	border Border
	trees  map[int]*Tree
	world  World
	exit   bool
}

// NOTE
// Interesting Unicode characters (e.g. arrows) start at 2190.
func (game *Game) DrawBorder(screen tcell.Screen) {
	for c := game.border.x1 + 1; c <= game.border.x2-1; c++ { // Add top and bottom borders
		screen.SetContent(c, game.border.y1, tcell.RuneHLine, nil, tcell.StyleDefault)
		screen.SetContent(c, game.border.y2, tcell.RuneHLine, nil, tcell.StyleDefault)
	}

	for r := game.border.y1 + 1; r <= game.border.y2-1; r++ { // Add left and right borders
		screen.SetContent(game.border.x1, r, tcell.RuneVLine, nil, tcell.StyleDefault)
		screen.SetContent(game.border.x2, r, tcell.RuneVLine, nil, tcell.StyleDefault)
	}

	// Add corners
	screen.SetContent(game.border.x1, game.border.y1, tcell.RuneULCorner, nil, tcell.StyleDefault)
	screen.SetContent(game.border.x2, game.border.y1, tcell.RuneURCorner, nil, tcell.StyleDefault)
	screen.SetContent(game.border.x1, game.border.y2, tcell.RuneLLCorner, nil, tcell.StyleDefault)
	screen.SetContent(game.border.x2, game.border.y2, tcell.RuneLRCorner, nil, tcell.StyleDefault)
}

func (game *Game) DrawTrees(screen tcell.Screen) {
	for _, tree := range game.trees {
		switch tree.state {
		case TreeStateStump:
			screen.SetContent(tree.x, tree.y, '▄', nil, tcell.StyleDefault)
		case TreeStateTrunk:
			screen.SetContent(tree.x, tree.y, '█', nil, tcell.StyleDefault)
		case TreeStateStumpling:
			screen.SetContent(tree.x, tree.y, '╻', nil, tcell.StyleDefault)
		case TreeStateSapling:
			screen.SetContent(tree.x, tree.y, '┃', nil, tcell.StyleDefault)
		case TreeStateSeed:
			screen.SetContent(tree.x, tree.y, '.', nil, tcell.StyleDefault)
		case TreeStateAdult:
			screen.SetContent(tree.x, tree.y, '█', nil, tcell.StyleDefault)
			screen.SetContent(tree.x, tree.y-1, '▄', nil, tcell.StyleDefault)
		}
	}
}

func (game *Game) ClearPlayer(screen tcell.Screen) {
	screen.SetContent(game.player.x, game.player.y, ' ', nil, tcell.StyleDefault)
}

func (game *Game) DrawPlayer(screen tcell.Screen) {
	w, h := screen.Size()
	screen.SetContent(w/2, h/2, '@', nil, tcell.StyleDefault) // Draw the player at the "center" of the view
}

// Only draw things within the player view range
func (game *Game) DrawViewport(screen tcell.Screen) {
	player_position := Coordinate{game.player.x, game.player.y}
	game.DrawPlayer(screen)

	/*
		view_radius = 5
		y_up = 5
		y_down = -5
		x_left = 5
		x_right = -5
				.....................
				.....................
				.....................
				.....................
				.....................
				..........@..........
				.....................
				.....................
				.....................
				.....................
				.....................
		Player pos = (5,5)
	*/

	for x := player_position.x - game.player.vision_radius; x <= player_position.x+game.player.vision_radius; x++ {
		for y := player_position.y - game.player.vision_radius; y <= player_position.y+game.player.vision_radius; y++ {
			obj, found := game.world.content[Coordinate{x, y}]
			if found {
				w, h := screen.Size()
				obj_viewport_x := (w / 2) + (x - game.player.x) // Player_viewport_x + Object_real_x - Player_real_x
				obj_viewport_y := (h / 2) + (y - game.player.y) // Player_viewport_y + Object_real_y - Player_real_y
				switch obj.(type) {
				case Object:
					// Draw object
					screen.SetContent(obj_viewport_x, obj_viewport_y, '#', nil, tcell.StyleDefault)
				case *Tree:
					tree := obj.(*Tree)
					// Draw tree
					switch tree.state {
					case TreeStateStump:
						screen.SetContent(obj_viewport_x, obj_viewport_y, '▄', nil, tcell.StyleDefault)
					case TreeStateTrunk:
						screen.SetContent(obj_viewport_x, obj_viewport_y, '█', nil, tcell.StyleDefault)
					case TreeStateStumpling:
						screen.SetContent(obj_viewport_x, obj_viewport_y, '╻', nil, tcell.StyleDefault)
					case TreeStateSapling:
						screen.SetContent(obj_viewport_x, obj_viewport_y, '┃', nil, tcell.StyleDefault)
					case TreeStateSeed:
						screen.SetContent(obj_viewport_x, obj_viewport_y, '.', nil, tcell.StyleDefault)
					case TreeStateAdult:
						screen.SetContent(obj_viewport_x, obj_viewport_y, '█', nil, tcell.StyleDefault)
						screen.SetContent(obj_viewport_x, obj_viewport_y-1, '▄', nil, tcell.StyleDefault)
					}
				}
			}
		}
	}
}

func (game *Game) Draw(screen tcell.Screen) {
	screen.Clear()
	game.DrawViewport(screen)
	screen.Show()
}

func (game *Game) MovePlayer(screen tcell.Screen, len int, dir int) {
	game.ClearPlayer(screen)

	// Determine (potential) new location.
	pMoved := game.player
	newX := game.player.x
	newY := game.player.y

	if len != 0 {
		switch dir {
		case DirUp:
			newY = game.player.y - len
		case DirRight:
			newX = game.player.x + len
		case DirDown:
			newY = game.player.y + len
		case DirLeft:
			newX = game.player.x - len
		}
	}

	// Prevent player from moving through an impassable object
	obj, exists := game.world.content[Coordinate{newX, newY}]
	if exists {
		switch obj.(type) {
		case Object, *Tree:
			newY = pMoved.y
			newX = pMoved.x
		default:
			pMoved.y = newY
			pMoved.x = newX
		}
	} else {
		pMoved.y = newY
		pMoved.x = newX
	}

	game.player = pMoved
}

func (game *Game) AddSeeds() int {
	seedCount := 0

	// Get max index of current trees map
	var maxIndex int
	for index := range game.trees {
		if maxIndex < index {
			maxIndex = index
		}
	}

	// Possibly create new seeds
	for i := 1; i <= SeedCreationMax; i++ {
		if rand.Float64() <= SeedCreationChance {
			var x, y int
			overlaps := false
			// Prevent seed from spawning on occupied point
			for {
				x = rand.Intn(game.border.x2-1) + game.border.x1
				y = rand.Intn(game.border.y2-1) + game.border.y1
				for _, tree := range game.trees {
					if x == tree.x && y == tree.y {
						overlaps = true
						break
					}
				}
				if !overlaps {
					break
				}

			}
			game.trees[maxIndex+i] = &Tree{x, y, TreeStateSeed}
			seedCount++
		}
	}

	return seedCount
}

func (game *Game) GrowTrees() int {
	growthCount := 0

	for _, object := range game.world.content {
		switch object.(type) {
		case *Tree:
			tree := object.(*Tree)
			if tree.state == TreeStateSeed && rand.Float64() <= GrowthChanceSeed {
				tree.state = TreeStateSapling
				growthCount++
			} else if tree.state == TreeStateSapling && rand.Float64() <= GrowthChanceSapling {
				tree.state = TreeStateAdult
				growthCount++
			}
		}
	}

	return growthCount
}

func (game *Game) PopulateTrees(screen tcell.Screen) int {
	// states := []int{
	// 	TreeStateSeed,
	// 	TreeStateSapling,
	// 	TreeStateAdult,
	// }
	maxTreeCount := rand.Intn(30) + 10
	treeCount := 0
	for i := 0; i < maxTreeCount; i++ {
		state := TreeStateSeed //states[rand.Intn(3)]
		x := rand.Intn(game.world.width)
		y := rand.Intn(game.world.height)
		game.world.content[Coordinate{x, y}] = &Tree{x, y, state}
		treeCount++
	}

	return treeCount
}

func (game *Game) DecrementTree(screen tcell.Screen, indexToRemove int) bool {
	// adult ------> trunk
	// trunk ------> stump
	// stump ------> removed
	// sapling ----> stumpling
	// stumpling --> removed
	tree, exists := game.trees[indexToRemove]

	if !exists {
		return false
	}

	switch tree.state {
	case TreeStateAdult:
		tree.state = TreeStateTrunk
		return true
	case TreeStateTrunk:
		tree.state = TreeStateStump
		return true
	case TreeStateSapling:
		tree.state = TreeStateStumpling
		return true
	case TreeStateStump, TreeStateStumpling:
		delete(game.trees, indexToRemove)
		return true
	}

	return false
}

func (game *Game) Chop(screen tcell.Screen, dir int) int {
	choppedCount := 0
outside:
	for index, tree := range game.trees {
		if tree.state != TreeStateRemoved && tree.state != TreeStateSeed {
			isAbove := tree.y == game.player.y-1 && tree.x == game.player.x
			isRight := tree.y == game.player.y && tree.x == game.player.x+1
			isBelow := tree.y == game.player.y+1 && tree.x == game.player.x
			isLeft := tree.y == game.player.y && tree.x == game.player.x-1
			switch dir {
			case DirOmni:
				if (isAbove || isRight || isBelow || isLeft) && game.DecrementTree(screen, index) {
					choppedCount++
				}
			case DirUp:
				if isAbove && game.DecrementTree(screen, index) {
					choppedCount++
					break outside
				}
			case DirRight:
				if isRight && game.DecrementTree(screen, index) {
					choppedCount++
					break outside
				}
			case DirDown:
				if isBelow && game.DecrementTree(screen, index) {
					choppedCount++
					break outside
				}
			case DirLeft:
				if isLeft && game.DecrementTree(screen, index) {
					choppedCount++
					break outside
				}
			}
		}
	}

	return choppedCount
}
