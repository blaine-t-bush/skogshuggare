package main

import "math/rand"

func (game *Game) SpreadFire() int {
	spreadCount := 0

	for position, content := range game.world.content {
		switch content := content.(type) {
		case Object:
			if content.category == ContentCategoryFire && rand.Float64() <= FireSpreadChance {
				// Pick random direction without fire
				game.AppendToMenuMessages("spread")
				spreadCoordinate := position
				switch GetRandomDirection() {
				case DirUp:
					spreadCoordinate.y = spreadCoordinate.y - 1
				case DirRight:
					spreadCoordinate.x = spreadCoordinate.x + 1
				case DirDown:
					spreadCoordinate.y = spreadCoordinate.y + 1
				case DirLeft:
					spreadCoordinate.x = spreadCoordinate.x - 1
				}

				// Check if blocked
				spread := true
				if existingContent, exists := game.world.content[spreadCoordinate]; exists {
					switch existingContent := existingContent.(type) {
					case Object:
						if existingContent.category == ContentCategoryTerrain {
							spread = false
						}
					}
				}

				if spread {
					game.world.content[spreadCoordinate] = Object{KeyFire, ContentCategoryFire, false, false}
					spreadCount++
				}
			}
		}
	}

	return spreadCount
}
