package main

import (
	"math/rand"
	"sync"
	"time"
)

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
}
