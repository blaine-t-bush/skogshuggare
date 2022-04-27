package main

import (
	"math/rand"
	"sync"
	"time"
)

func (game *Game) AnimationHandler(wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()

	for {
		// Stop animating if game is closed.
		if game.exit {
			return
		}

		mutex.Lock()

		// Randomly change fire and water glyphs according to available keys.
		fireKeys := [2]int{KeyFireLight, KeyFireHeavy}
		waterKeys := [2]int{KeyWaterLight, KeyWaterHeavy}
		for _, content := range game.world.content {
			switch content := content.(type) {
			case *Object:
				if content.key == KeyWaterLight || content.key == KeyWaterHeavy {
					content.key = waterKeys[rand.Intn(len(waterKeys))]
				}
			case *Fire:
				content.key = fireKeys[rand.Intn(len(fireKeys))]
			}
		}

		// Re-render with the updated glyphs.
		game.Draw()

		mutex.Unlock()

		// Wait AnimationRate milliseconds before updating animation states.
		time.Sleep(AnimationRate * time.Millisecond)
	}
}
