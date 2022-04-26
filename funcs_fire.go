package main

import "math/rand"

func (game *Game) SpreadFire() int {
	spreadCount := 0

	for position, content := range game.world.content {
		switch content := content.(type) {
		case Object:
			if content.category == ContentCategoryFire && rand.Float64() <= FireSpreadChance {
				// Pick random direction
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

				// Spread if not blocked
				if spread {
					game.world.content[spreadCoordinate] = Object{KeyFire, ContentCategoryFire, false, false}
					spreadCount++
				}
			}
		}
	}

	return spreadCount
}

func (game *Game) CheckFireDamage() int {
	damage := 0

	// Check if fire exists on player tile
	if content, exists := game.world.content[game.player.position]; exists {
		switch content := content.(type) {
		case Object:
			if content.category == ContentCategoryFire {
				// Damage player
				newHitPoints := game.player.hitPointsCurrent - DamageFire // FIXME check if player HP reached 0
				game.player.hitPointsCurrent = newHitPoints
				damage++
			}
		}
	}

	// Check if fire exists on squirrel tiles
	for key, squirrel := range game.squirrels {
		if content, exists := game.world.content[squirrel.position]; exists {
			switch content := content.(type) {
			case Object:
				if content.category == ContentCategoryFire {
					// Damage squirrel
					newHitPoints := squirrel.hitPointsCurrent - DamageFire
					if newHitPoints <= 0 {
						// Delete squirrel
						delete(game.squirrels, key)
					} else {
						squirrel.hitPointsCurrent = newHitPoints
					}
					damage++
				}
			}
		}
	}

	return damage
}
