package main

import (
	"math/rand"

	"github.com/gdamore/tcell"
)

type Actor struct {
	x int // Player x-coordinate
	y int // Player y-coordinate
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

const (
	// Game parameters
	TickRate = 30 // Milliseconds between ticks
	// Actors
	ActorPlayer   = 1
	ActorSquirrel = 2
	// Characters
	CharacterPlayer   = '@'
	CharacterSquirrel = '~'
	// Directions
	DirUp     = 0
	DirRight  = 1
	DirDown   = 2
	DirLeft   = 3
	DirOmni   = 4
	DirRandom = 5
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

type Game struct {
	player   Actor
	squirrel Actor
	border   Border
	trees    map[int]*Tree
	exit     bool
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

func (game *Game) ClearActor(screen tcell.Screen, actorType int) {
	var actor Actor
	switch actorType {
	case ActorPlayer:
		actor = game.player
	case ActorSquirrel:
		actor = game.squirrel
	}
	screen.SetContent(actor.x, actor.y, ' ', nil, tcell.StyleDefault)
}

func (game *Game) DrawActor(screen tcell.Screen, actorType int) {
	var actor Actor
	var character rune
	switch actorType {
	case ActorPlayer:
		actor = game.player
		character = CharacterPlayer
	case ActorSquirrel:
		actor = game.squirrel
		character = CharacterSquirrel
	}
	screen.SetContent(actor.x, actor.y, character, nil, tcell.StyleDefault)
}

func (game *Game) Draw(screen tcell.Screen) {
	screen.Clear()
	game.DrawBorder(screen)
	game.DrawActor(screen, ActorPlayer)
	game.DrawActor(screen, ActorSquirrel)
	game.DrawTrees(screen)
	screen.Show()
}

func (game *Game) MoveActor(screen tcell.Screen, actorType int, len int, dir int) {
	var actor Actor
	switch actorType {
	case ActorPlayer:
		actor = game.player
	case ActorSquirrel:
		actor = game.squirrel
	}

	game.ClearActor(screen, actorType)

	// Determine (potential) new location.
	if dir == DirRandom {
		randInt := rand.Intn(4)
		switch randInt {
		case 0:
			dir = DirUp
		case 1:
			dir = DirRight
		case 2:
			dir = DirDown
		case 3:
			dir = DirLeft
		}
	}
	aMoved := actor
	if len != 0 {
		switch dir {
		case DirUp:
			aMoved.y = actor.y - len
		case DirRight:
			aMoved.x = actor.x + len
		case DirDown:
			aMoved.y = actor.y + len
		case DirLeft:
			aMoved.x = actor.x - len
		}
	}

	// Prevent update if player would collide with tree trunks, but not with
	// tree canopies.
	for _, tree := range game.trees {
		if tree.state != TreeStateRemoved && aMoved.y == tree.y && aMoved.x == tree.x {
			aMoved.x = actor.x
			aMoved.y = actor.y
		}
	}

	// Prevent update if new location is past left or right boundaries.
	if aMoved.x <= game.border.x1 {
		aMoved.x = game.border.x1 + game.border.t
	} else if aMoved.x >= game.border.x2 {
		aMoved.x = game.border.x2 - game.border.t
	}

	// Prevent update if new location is past top or bottom boundaries.
	if aMoved.y <= game.border.y1 {
		aMoved.y = game.border.y1 + game.border.t
	} else if aMoved.y >= game.border.y2 {
		aMoved.y = game.border.y2 - game.border.t
	}

	switch actorType {
	case ActorPlayer:
		game.player = aMoved
	case ActorSquirrel:
		game.squirrel = aMoved
	}
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
	for index, tree := range game.trees {
		if tree.state == TreeStateSeed && rand.Float64() <= GrowthChanceSeed {
			game.trees[index].state = TreeStateSapling
			growthCount++
		} else if tree.state == TreeStateSapling && rand.Float64() <= GrowthChanceSapling {
			game.trees[index].state = TreeStateAdult
			growthCount++
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
		x := rand.Intn(game.border.x2-1) + game.border.x1
		y := rand.Intn(game.border.y2-1) + game.border.y1
		game.trees[i] = &Tree{x, y, state}
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
