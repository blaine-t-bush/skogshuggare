package main

import (
	"math/rand"
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

func (game *Game) PopulateTrees() int {
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

func (game *Game) PopulateGrass() int {
	keys := []int{
		KeyGrassLight,
		KeyGrassHeavy,
	}
	maxGrassCount := rand.Intn(10) + 6
	grassCount := 0
	for i := 0; i < maxGrassCount; i++ {
		key := keys[rand.Intn(len(keys))]
		coordinate := game.GetRandomPlantableCoordinate()
		game.world.content[coordinate] = &AnimatedObject{key, GetRandomAnimationStage(key), false, true, false}
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

func (game *Game) DecrementTree(position Coordinate, stages int) bool {
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

func (game *Game) Chop(dir int, stages int) int {
	// Determine which coordinates to check for chopping based on direction and player position.
	var targetCoordinates [4]Coordinate
	switch dir {
	case DirOmni:
		targetCoordinates[0] = Translate(game.player.position, 0, -1)
		targetCoordinates[1] = Translate(game.player.position, 1, 0)
		targetCoordinates[2] = Translate(game.player.position, 0, 1)
		targetCoordinates[3] = Translate(game.player.position, -1, 0)
	case DirUp:
		targetCoordinates[0] = Translate(game.player.position, 0, -1)
	case DirRight:
		targetCoordinates[0] = Translate(game.player.position, 1, 0)
	case DirDown:
		targetCoordinates[0] = Translate(game.player.position, 0, 1)
	case DirLeft:
		targetCoordinates[0] = Translate(game.player.position, -1, 0)
	}

	// Chop trees that are within the target coordinate(s)
	choppedCount := 0
	for _, targetCoordinate := range targetCoordinates {
		if content, exists := game.world.content[targetCoordinate]; exists {
			switch content.(type) {
			case *Tree:
				if game.DecrementTree(targetCoordinate, stages) {
					choppedCount++
				}
			}
		}
	}

	return choppedCount
}
