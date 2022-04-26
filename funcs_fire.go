package main

import (
	"math/rand"
)

func BurnoutChance(t int) float64 {
	return 1 - (1 / (1 + float64(t)/float64(FireBurnoutHalflife)))
}

func (game *Game) UpdateFire() int {
	spreadCount := 0

	for position, content := range game.world.content {
		switch content := content.(type) {
		case *Fire:
			// Check for burnout
			if rand.Float64() <= BurnoutChance(content.age) {
				delete(game.world.content, position)
				game.world.content[position] = Object{KeyBurnt, false, false, false}
			}

			// Check for spreading
			if rand.Float64() <= FireSpreadChance {
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
						if !existingContent.flammable {
							spread = false
						}
					}
				}

				// Spread if not blocked
				if spread {
					game.world.content[spreadCoordinate] = &Fire{spreadCoordinate, 0}
					spreadCount++
				}
			}

			// Increment age
			content.age = content.age + 1
		}
	}

	return spreadCount
}

func (game *Game) CheckFireDamage() int {
	damage := 0

	// Check if fire exists on player tile
	if content, exists := game.world.content[game.player.position]; exists {
		switch content.(type) {
		case *Fire:
			// Damage player
			newHitPoints := game.player.hitPointsCurrent - DamageFire
			if newHitPoints <= 0 {
				game.exit = true
			}
			game.player.hitPointsCurrent = newHitPoints
			damage++
		}
	}

	// Check if fire exists on squirrel tiles
	for key, squirrel := range game.squirrels {
		if content, exists := game.world.content[squirrel.position]; exists {
			switch content.(type) {
			case *Fire:
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

	return damage
}
