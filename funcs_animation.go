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

func GetRandomAnimationStage(animationStatesKey int) int {
	return rand.Intn(len(animationStates[animationStatesKey]))
}

func GetNextAnimationStage(animationStatesKey int, currentStage int) int {
	contentAnimationStates := animationStates[animationStatesKey]
	var nextStage int
	if currentStage == len(contentAnimationStates)-1 {
		nextStage = 0
	} else {
		nextStage = currentStage + 1
	}
	return nextStage
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
