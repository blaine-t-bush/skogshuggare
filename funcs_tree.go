package main

import (
	"math/rand"

	"github.com/gdamore/tcell"
)

// func (game *Game) PlantSeed(x int, y int) bool {
// 	// Get max index of current trees map
// 	var maxIndex int
// 	for index := range game.trees {
// 		if maxIndex < index {
// 			maxIndex = index
// 		}
// 	}

// 	// Check if seed overlaps with existing trees
// 	// If it does not overlap, plant it
// 	if !game.PointIsBlocked(Coordinate{x, y}, game.GetNonActorBlockedPoints()) {
// 		game.trees[maxIndex+1] = &Tree{x, y, TreeStateSeed}
// 		return true
// 	} else {
// 		return false
// 	}
// }

func (game *Game) PopulateTrees(screen tcell.Screen) int {
	states := []int{
		TreeStateSeed,
		TreeStateSapling,
		TreeStateAdult,
	}
	maxTreeCount := rand.Intn(30) + 10
	treeCount := 0
	for i := 0; i < maxTreeCount; i++ {
		state := states[rand.Intn(len(states))]
		x := rand.Intn(game.world.width)
		y := rand.Intn(game.world.height)
		for {
			_, exists := game.world.content[Coordinate{x, y}]
			if !exists {
				game.world.content[Coordinate{x, y}] = &Tree{Coordinate{x, y}, state}
				break
			} else {
				x = rand.Intn(game.world.width)
				y = rand.Intn(game.world.height)
			}

		}
		treeCount++
	}

	return treeCount
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

func (game *Game) DecrementTree(screen tcell.Screen, position Coordinate) bool {
	// adult ------> trunk
	// trunk ------> stump
	// stump ------> removed
	// sapling ----> stumpling
	// stumpling --> removed
	object, exists := game.world.content[position]

	if !exists {
		return false
	}

	switch object.(type) {
	case *Tree:
		tree := object.(*Tree)
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
	for _, object := range game.world.content {
		switch object.(type) {
		case *Tree:
			tree := object.(*Tree)
			isAbove := tree.position.y == game.player.position.y-1 && tree.position.x == game.player.position.x
			isRight := tree.position.y == game.player.position.y && tree.position.x == game.player.position.x+1
			isBelow := tree.position.y == game.player.position.y+1 && tree.position.x == game.player.position.x
			isLeft := tree.position.y == game.player.position.y && tree.position.x == game.player.position.x-1
			switch dir {
			case DirOmni:
				if (isAbove || isRight || isBelow || isLeft) && game.DecrementTree(screen, tree.position) {
					choppedCount++
				}
			case DirUp:
				if isAbove && game.DecrementTree(screen, tree.position) {
					choppedCount++
					break outside
				}
			case DirRight:
				if isRight && game.DecrementTree(screen, tree.position) {
					choppedCount++
					break outside
				}
			case DirDown:
				if isBelow && game.DecrementTree(screen, tree.position) {
					choppedCount++
					break outside
				}
			case DirLeft:
				if isLeft && game.DecrementTree(screen, tree.position) {
					choppedCount++
					break outside
				}
			}
		}
	}

	return choppedCount
}
