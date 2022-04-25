package main

import (
	"github.com/gdamore/tcell"
)

func (game *Game) MoveSquirrel(screen tcell.Screen, len int, dir int, squirrelKey int) {
	game.MoveActor(screen, game.squirrels[squirrelKey], len, dir)
}

func (game *Game) MovePlayer(screen tcell.Screen, len int, dir int) {
	game.MoveActor(screen, &game.player, len, dir)
}

func (game *Game) MoveActor(screen tcell.Screen, actor *Actor, len int, dir int) {

	// Determine (potential) new location.
	if dir == DirRandom {
		dir = GetRandomDirection()
	}

	newX := actor.position.x
	newY := actor.position.y

	if len != 0 {
		switch dir {
		case DirUp:
			newY = actor.position.y - len
		case DirRight:
			newX = actor.position.x + len
		case DirDown:
			newY = actor.position.y + len
		case DirLeft:
			newX = actor.position.x - len
		}
	}

	// Prevent actor from moving through an collidable object.
	if content, exists := game.world.content[Coordinate{newX, newY}]; exists {
		switch content := content.(type) {
		case Object:
			if content.collidable {
				break
			} else {
				actor.position.y = newY
				actor.position.x = newX
			}
		case *Tree:
			break // Collide; do not change position.
		default:
			actor.position.y = newY
			actor.position.x = newX
		}
	} else {
		actor.position.y = newY
		actor.position.x = newX
	}
}

func (actor *Actor) IsAdjacentToDestination() bool {
	dest := actor.destination
	pos := actor.position
	if dest.y == pos.y && (dest.x == pos.x+1 || dest.x == pos.x-1) {
		return true
	} else if dest.x == pos.x && (dest.y == pos.y+1 || dest.y == pos.y-1) {
		return true
	}

	return false
}
