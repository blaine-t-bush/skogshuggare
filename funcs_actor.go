package main

import (
	"math/rand"

	"github.com/gdamore/tcell"
)

func GetRandomDirection() int {
	randInt := rand.Intn(4)
	switch randInt {
	case 0:
		return DirUp
	case 1:
		return DirRight
	case 2:
		return DirDown
	case 3:
		return DirLeft
	default:
		return DirNone
	}
}

func (game *Game) MoveActor(screen tcell.Screen, actorType int, len int, dir int) {
	// Determine which actor to update.
	var actor *Actor
	switch actorType {
	case ActorPlayer:
		actor = &game.player
	case ActorSquirrel:
		actor = &game.squirrel
	}

	// Un-draw the actor at its existing location.
	game.ClearActor(screen, actorType)

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
	obj, exists := game.world.content[Coordinate{newX, newY}]
	if exists {
		switch obj.(type) {
		case Object, *Tree:
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
