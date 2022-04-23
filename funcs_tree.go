package main

import (
	"math/rand"

	"github.com/gdamore/tcell"
)

func (game *Game) PlantSeed(coordinate Coordinate) bool {
	// Get max index of current trees map
	if !game.IsBlocked(coordinate) {
		game.world.content[coordinate] = &Tree{coordinate, TreeStateSeed}
		return true
	} else {
		return false
	}
}

func (game *Game) PopulateTrees(screen tcell.Screen) int {
	states := []int{
		TreeStateSeed,
		TreeStateSapling,
		TreeStateAdult,
	}
	maxTreeCount := rand.Intn(5) + 3
	treeCount := 0
	for i := 0; i < maxTreeCount; i++ {
		state := states[rand.Intn(len(states))]
		coordinate := game.GetRandomAvailableCoordinate()
		game.world.content[coordinate] = &Tree{coordinate, state}
		treeCount++
	}

	return treeCount
}

func (game *Game) PopulateGrass(screen tcell.Screen) int {
	maxGrassCount := rand.Intn(10) + 6
	grassCount := 0
	for i := 0; i < maxGrassCount; i++ {
		coordinate := game.GetRandomAvailableCoordinate()
		if rand.Intn(2) == 0 {
			game.world.content[coordinate] = Object{'"', false}
		} else {
			game.world.content[coordinate] = Object{'\'', false}
		}
		grassCount++
	}

	return grassCount
}

func (game *Game) GrowTrees() int {
	growthCount := 0

	for _, content := range game.world.content {
		switch content := content.(type) {
		case *Tree:
			if content.state == TreeStateSeed && rand.Float64() <= GrowthChanceSeed {
				content.state = TreeStateSapling
				growthCount++
			} else if content.state == TreeStateSapling && rand.Float64() <= GrowthChanceSapling {
				content.state = TreeStateAdult
				growthCount++
			}
		}
	}

	return growthCount
}

func (game *Game) DecrementTree(screen tcell.Screen, position Coordinate) bool {
	// adult ------> trunk
	// trunk ------> stump
	// stump ------> removed
	// sapling ----> stumpling
	// stumpling --> removed
	content, exists := game.world.content[position]

	if !exists {
		return false
	}

	switch content := content.(type) {
	case *Tree:
		switch content.state {
		case TreeStateAdult:
			content.state = TreeStateTrunk
			return true
		case TreeStateTrunk:
			content.state = TreeStateStump
			return true
		case TreeStateSapling:
			content.state = TreeStateStumpling
			return true
		case TreeStateStump, TreeStateStumpling:
			delete(game.world.content, position)
			return true
		}
	default:
		return false // Coordinate does not correspond to a tree
	}

	return false
}

func (game *Game) Chop(screen tcell.Screen, dir int) int {
	choppedCount := 0
outside:
	for _, content := range game.world.content {
		switch content := content.(type) {
		case *Tree:
			isAbove := content.position.y == game.player.position.y-1 && content.position.x == game.player.position.x
			isRight := content.position.y == game.player.position.y && content.position.x == game.player.position.x+1
			isBelow := content.position.y == game.player.position.y+1 && content.position.x == game.player.position.x
			isLeft := content.position.y == game.player.position.y && content.position.x == game.player.position.x-1
			switch dir {
			case DirOmni:
				if (isAbove || isRight || isBelow || isLeft) && game.DecrementTree(screen, content.position) {
					choppedCount++
				}
			case DirUp:
				if isAbove && game.DecrementTree(screen, content.position) {
					choppedCount++
					break outside
				}
			case DirRight:
				if isRight && game.DecrementTree(screen, content.position) {
					choppedCount++
					break outside
				}
			case DirDown:
				if isBelow && game.DecrementTree(screen, content.position) {
					choppedCount++
					break outside
				}
			case DirLeft:
				if isLeft && game.DecrementTree(screen, content.position) {
					choppedCount++
					break outside
				}
			}
		}
	}

	return choppedCount
}
