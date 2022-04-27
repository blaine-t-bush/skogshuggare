package main

import (
	"math/rand"

	"github.com/gdamore/tcell"
)

func BurnoutChance(t int) float64 {
	return 1 - (1 / (1 + float64(t)/float64(FireBurnoutHalflife)))
}

func (game *Game) SpawnRandomFire() {
	coord := game.GetRandomFlammableCoordinate()
	game.world.content[coord] = &Fire{coord, 0}
}

func (game *Game) UpdateFire() int {
	spreadAndSpawnCount := 0

	// Check for spreading and burning out of existing fire.
	for position, content := range game.world.content {
		switch content := content.(type) {
		case *Fire:
			// Check for burnout
			if rand.Float64() <= BurnoutChance(content.age) {
				delete(game.world.content, position)
				game.world.content[position] = Object{KeyBurnt, false, false, true}
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
					spreadAndSpawnCount++
				}
			}

			// Increment age
			content.age = content.age + 1
		}
	}

	// Check for spawning of new fires
	if rand.Float64() <= FireSpawnChance {
		game.SpawnRandomFire()
	}

	return spreadAndSpawnCount
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

func (game *Game) Dig(screen tcell.Screen, dir int) int {
	// Determine which coordinates to check for digging based on direction and player position.
	var targetCoordinates [4]Coordinate
	switch dir {
	case DirOmni:
		targetCoordinates[0] = Coordinate{game.player.position.x, game.player.position.y - 1}
		targetCoordinates[1] = Coordinate{game.player.position.x + 1, game.player.position.y}
		targetCoordinates[2] = Coordinate{game.player.position.x, game.player.position.y + 1}
		targetCoordinates[3] = Coordinate{game.player.position.x - 1, game.player.position.y}
	case DirUp:
		targetCoordinates[0] = Coordinate{game.player.position.x, game.player.position.y - 1}
	case DirRight:
		targetCoordinates[0] = Coordinate{game.player.position.x + 1, game.player.position.y}
	case DirDown:
		targetCoordinates[0] = Coordinate{game.player.position.x, game.player.position.y + 1}
	case DirLeft:
		targetCoordinates[0] = Coordinate{game.player.position.x - 1, game.player.position.y}
	}

	// Dig tiles that are within the target coordinate(s) and unblocked
	dugCount := 0
	for _, targetCoordinate := range targetCoordinates {
		dig := true
		if content, exists := game.world.content[targetCoordinate]; exists {
			switch content := content.(type) {
			case Object:
				if content.collidable {
					dig = false
				}
			case *Tree:
				dig = false
			}
		}

		if dig {
			game.world.content[targetCoordinate] = Object{KeyFirebreak, false, false, false}
			dugCount++
		}
	}
	return dugCount
}
