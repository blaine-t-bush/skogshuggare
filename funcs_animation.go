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
	for range ticker.C {
		mutex.Lock()
		game.AnimationUpdate()
		game.Draw()
		mutex.Unlock()
		if game.exit {
			wg.Done()
			return
		}
	}

}

func (game *Game) AnimationUpdate() {
	// Randomly change fire and water glyphs according to available keys.
	for _, content := range game.world.content {
		switch content := content.(type) {
		case *AnimatedObject:
			content.animationStage = GetNextAnimationStage(content.key, content.animationStage)
		case *Fire:
			content.animationStage = GetNextAnimationStage(KeyFire, content.animationStage)
		}
	}
}
