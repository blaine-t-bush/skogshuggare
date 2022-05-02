package main

import (
	"math/rand"
	"sync"
	"time"
)

func IsAnimatedObject(content interface{}) bool {
	isAnimatedObject := false
	switch content.(type) {
	case *AnimatedObject:
		isAnimatedObject = true
	}
	return isAnimatedObject
}

func GetAnimationState(key int, stage int) int {
	return animationMarkov[key][stage].state
}

func GetRandomAnimationStage(key int) int {
	return rand.Intn(len(animationMarkov[key]))
}

func GetNextAnimationStage(key int, stage int) int {
	// Get Markov connections for current stage
	connections := animationMarkov[key][stage].connections
	// Get total probability for normalization, in case it doesn't sum to 1
	totalProbability := 0.0
	for _, connection := range connections {
		totalProbability = totalProbability + connection.probability
	}
	// Create normalized probability ranges
	cumulativeProbability := 0.0
	probabilityRanges := make(map[int]float64)
	for _, connection := range connections {
		probabilityRanges[connection.stage] = (connection.probability + cumulativeProbability) / totalProbability
		cumulativeProbability = cumulativeProbability + connection.probability
	}

	// Choose state
	roll := rand.Float64()
	var newStage int
	for stage, probabilityRange := range probabilityRanges {
		if roll <= probabilityRange {
			newStage = stage
			break
		}
	}

	return newStage
}

func (game *Game) AnimationTicker(wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()

	// Initialize animation update ticker.
	ticker := time.NewTicker(AnimationTickDuration * time.Millisecond)

	// Update animation state and re-draw on every tick.
	counter := 0
	for range ticker.C {
		counter++
		mutex.Lock()
		game.AnimationUpdate(counter)
		game.Draw()
		mutex.Unlock()
		if game.exit {
			wg.Done()
			return
		}
	}

}

func (game *Game) AnimationUpdate(counter int) {
	// Check for chance to spawn birds and clouds
	if rand.Float64() <= BirdSpawnChance {
		game.world.decorations[game.GetRandomBorderCoordinate()] = &DecorationObject{KeyBird, GetRandomAnimationStage(KeyBird), GetRandomDirection()}
	}

	// Randomly change fire and water glyphs according to available keys.
	for _, content := range game.world.content {
		switch content := content.(type) {
		case *AnimatedObject:
			if counter%animationRates[content.key] == 0 {
				content.animationStage = GetNextAnimationStage(content.key, content.animationStage)
			}
		case *Fire:
			if counter%animationRates[KeyFire] == 0 {
				content.animationStage = GetNextAnimationStage(KeyFire, content.animationStage)
			}
		}
	}

	// Update decoration animation stages and positions
	newDecorations := make(map[Coordinate]interface{})
	var oldCoords, newCoords []Coordinate
	for coord, decoration := range game.world.decorations {
		switch decoration := decoration.(type) {
		case *DecorationObject:
			if counter%animationRates[decoration.key] == 0 {
				newCoord := TranslateByDir(coord, decoration.direction, 1)
				newCoords = append(newCoords, newCoord)
				oldCoords = append(oldCoords, coord)
				newDecorations[newCoord] = &DecorationObject{
					key:            decoration.key,
					animationStage: GetNextAnimationStage(decoration.key, decoration.animationStage),
					direction:      decoration.direction,
				}
			}
		}
	}

	// Update decorations
	for coord, decoration := range newDecorations {
		switch decoration := decoration.(type) {
		case *DecorationObject:
			game.world.decorations[coord] = decoration
		}
	}

	// Determine which elements to delete from decorations map
	for _, oldCoord := range oldCoords {
		deleteCoord := true
		for _, newCoord := range newCoords {
			if oldCoord == newCoord {
				deleteCoord = false
			}
		}

		if deleteCoord {
			delete(game.world.decorations, oldCoord)
		}
	}
}
