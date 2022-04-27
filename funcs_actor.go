package main

func (game *Game) MoveSquirrel(len int, dir int, squirrelKey int) {
	game.MoveActor(game.squirrels[squirrelKey], len, dir)
}

func (game *Game) MovePlayer(len int, dir int) {
	game.MoveActor(&game.player, len, dir)
}

func (game *Game) MoveActor(actor *Actor, len int, dir int) bool {

	// Determine (potential) new location.
	if dir == DirRandom {
		dir = GetRandomDirection()
	}

	deltaX := 0
	deltaY := 0

	if len != 0 {
		switch dir {
		case DirUp:
			deltaY = -len
		case DirRight:
			deltaX = len
		case DirDown:
			deltaY = len
		case DirLeft:
			deltaX = -len
		}
	}

	// Prevent actor from moving through an collidable object.
	translate := true
	if content, exists := game.world.content[Translate(actor.position, deltaX, deltaY)]; exists {
		switch content := content.(type) {
		case Object:
			if content.collidable {
				translate = false
			}
		case *Tree:
			translate = false
		}
	}

	if translate {
		actor.position.Translate(deltaX, deltaY)
		return true
	}

	return false
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
