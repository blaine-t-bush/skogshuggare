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
		coordinate := game.GetRandomPlantableCoordinate()
		game.world.content[coordinate] = &Tree{coordinate, state}
		treeCount++
	}

	return treeCount
}

func (game *Game) PopulateGrass(screen tcell.Screen) int {
	keys := []int{
		KeyGrassLight,
		KeyGrassHeavy,
	}
	maxGrassCount := rand.Intn(10) + 6
	grassCount := 0
	for i := 0; i < maxGrassCount; i++ {
		key := keys[rand.Intn(len(keys))]
		coordinate := game.GetRandomPlantableCoordinate()
		game.world.content[coordinate] = Object{key, ContentCategoryDecoration, false, false}
		grassCount++
	}

	return grassCount
}

func (game *Game) GrowTrees() int {
	growthCount := 0

	for _, content := range game.world.content {
		switch content := content.(type) {
		case *Tree:
			if growthInfo, exists := treeGrowingStages[content.state]; exists {
				if rand.Float64() <= growthInfo.chance {
					content.state = growthInfo.newState
					growthCount++
				}
			}
		}
	}

	return growthCount
}

func (game *Game) DecrementTree(screen tcell.Screen, position Coordinate, stages int) bool {
	content, exists := game.world.content[position]

	if !exists { // No content of any type at this location
		return false
	}

	switch content := content.(type) {
	case *Tree:
		var exists bool
		newState := content.state
		for i := 0; i < stages; i++ { // We move down the harvesting stages one or more times
			newState, exists = treeHarvestingStages[newState]
			if !exists {
				return false
			}
		}

		if newState == TreeStateRemoved {
			delete(game.world.content, position)
			game.player.score++ // Increase player score when tree is felled
		}

		content.state = newState
		return true
	default:
		return false // Coordinate does not correspond to a tree
	}
}

func (game *Game) Chop(screen tcell.Screen, dir int, stages int) int {
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
				if (isAbove || isRight || isBelow || isLeft) && game.DecrementTree(screen, content.position, stages) {
					choppedCount++
				}
			case DirUp:
				if isAbove && game.DecrementTree(screen, content.position, stages) {
					choppedCount++
					break outside
				}
			case DirRight:
				if isRight && game.DecrementTree(screen, content.position, stages) {
					choppedCount++
					break outside
				}
			case DirDown:
				if isBelow && game.DecrementTree(screen, content.position, stages) {
					choppedCount++
					break outside
				}
			case DirLeft:
				if isLeft && game.DecrementTree(screen, content.position, stages) {
					choppedCount++
					break outside
				}
			}
		}
	}

	return choppedCount
}
