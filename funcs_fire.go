package main

import (
	"math/rand"
)

func IsFire(content interface{}) bool {
	isFire := false
	switch content := content.(type) {
	case *AnimatedObject:
		if content.key == KeyFire {
			isFire = true
		}
	}
	return isFire
}

func BurnoutChance(t int) float64 {
	return 1 - (1 / (1 + float64(t)/float64(FireBurnoutHalflife)))
}

func (game *Game) CreateFire(coordinate Coordinate) bool {
	game.world.content[coordinate] = &Fire{GetRandomAnimationStage(KeyFire), coordinate, 0}
	return true
}

func (game *Game) SpawnRandomFire() bool {
	return game.CreateFire(game.GetRandomFlammableCoordinate())
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
				game.world.content[position] = &StaticObject{KeyBurnt, false, false, true}
			}

			// Check for spreading
			if rand.Float64() <= FireSpreadChance {
				// Pick random direction
				deltaX := 0
				deltaY := 0
				switch GetRandomDirection() {
				case DirUp:
					deltaY = -1
				case DirRight:
					deltaX = 1
				case DirDown:
					deltaY = 1
				case DirLeft:
					deltaX = -1
				}

				// Check if blocked
				spread := true
				spreadCoordinate := Translate(position, deltaX, deltaY)
				if existingContent, exists := game.world.content[spreadCoordinate]; exists {
					switch existingContent := existingContent.(type) {
					case *StaticObject:
						if !existingContent.flammable {
							spread = false
						}
					}
				}

				// Spread if not blocked
				if spread {
					if game.CreateFire(spreadCoordinate) {
						spreadAndSpawnCount++
					}
				}
			}

			// Increment age
			content.age = content.age + 1
		}
	}

	// Check for spawning of new fires
	if rand.Float64() <= FireSpawnChance {
		if game.SpawnRandomFire() {
			spreadAndSpawnCount++
		}
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

func (game *Game) Dig(dir int) int {
	// Determine which coordinates to check for digging based on direction and player position.
	var targetCoordinates [4]Coordinate
	switch dir {
	case DirOmni:
		targetCoordinates[0] = Translate(game.player.position, 0, -1)
		targetCoordinates[1] = Translate(game.player.position, 1, 0)
		targetCoordinates[2] = Translate(game.player.position, 0, 1)
		targetCoordinates[3] = Translate(game.player.position, -1, 0)
	case DirUp:
		targetCoordinates[0] = Translate(game.player.position, 0, -1)
	case DirRight:
		targetCoordinates[0] = Translate(game.player.position, 1, 0)
	case DirDown:
		targetCoordinates[0] = Translate(game.player.position, 0, 1)
	case DirLeft:
		targetCoordinates[0] = Translate(game.player.position, -1, 0)
	}

	// Dig tiles that are within the target coordinate(s) and unblocked
	dugCount := 0
	for _, targetCoordinate := range targetCoordinates {
		dig := true
		if content, exists := game.world.content[targetCoordinate]; exists {
			switch content := content.(type) {
			case *StaticObject:
				if content.collidable {
					dig = false
				}
			case *Tree:
				dig = false
			}
		}

		if dig {
			game.world.content[targetCoordinate] = &StaticObject{KeyFirebreak, false, false, false}
			dugCount++
		}
	}
	return dugCount
}
