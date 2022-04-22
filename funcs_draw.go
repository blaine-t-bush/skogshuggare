package main

import "github.com/gdamore/tcell"

func (game *Game) Draw(screen tcell.Screen) {
	screen.Clear()
	game.DrawViewport(screen)
	screen.Show()
}

// Only draw things within the player view range
func (game *Game) DrawViewport(screen tcell.Screen) {
	game.DrawPlayer(screen)
	game.DrawSquirrel(screen)

	/*
		view_radius = 5
		y_up = 5
		y_down = -5
		x_left = 5
		x_right = -5
				.....................
				.....................
				.....................
				.....................
				.....................
				..........@..........
				.....................
				.....................
				.....................
				.....................
				.....................
		Player pos = (5,5)
	*/

	for x := game.player.position.x - game.player.visionRadius; x <= game.player.position.x+game.player.visionRadius; x++ {
		for y := game.player.position.y - game.player.visionRadius; y <= game.player.position.y+game.player.visionRadius; y++ {
			coord := Coordinate{x, y}

			w, h := screen.Size()
			// Get the viewport coordinates
			objViewportX := (w / 2) + (x - game.player.position.x) // Player_viewport_x + Object_real_x - Player_real_x
			objViewportY := (h / 2) + (y - game.player.position.y) // Player_viewport_y + Object_real_y - Player_real_y

			border, isBorder := game.world.borders[coord]
			if isBorder {
				switch border {
				case TopBorder, BottomBorder:
					screen.SetContent(objViewportX, objViewportY, tcell.RuneHLine, nil, tcell.StyleDefault)
				case RightBorder, LeftBorder:
					screen.SetContent(objViewportX, objViewportY, tcell.RuneVLine, nil, tcell.StyleDefault)
				case TopLeftCorner:
					screen.SetContent(objViewportX, objViewportY, tcell.RuneULCorner, nil, tcell.StyleDefault)
				case TopRightCorner:
					screen.SetContent(objViewportX, objViewportY, tcell.RuneURCorner, nil, tcell.StyleDefault)
				case BottomRightCorner:
					screen.SetContent(objViewportX, objViewportY, tcell.RuneLRCorner, nil, tcell.StyleDefault)
				case BottomLeftCorner:
					screen.SetContent(objViewportX, objViewportY, tcell.RuneLLCorner, nil, tcell.StyleDefault)
				}
				continue
			}

			obj, found := game.world.content[coord]
			if found {
				switch obj.(type) {
				case Object:
					// Draw object
					screen.SetContent(objViewportX, objViewportY, '#', nil, tcell.StyleDefault)
				case *Tree:
					tree := obj.(*Tree)
					// Draw tree
					switch tree.state {
					case TreeStateStump:
						screen.SetContent(objViewportX, objViewportY, '▄', nil, tcell.StyleDefault)
					case TreeStateTrunk:
						screen.SetContent(objViewportX, objViewportY, '█', nil, tcell.StyleDefault)
					case TreeStateStumpling:
						screen.SetContent(objViewportX, objViewportY, '╻', nil, tcell.StyleDefault)
					case TreeStateSapling:
						screen.SetContent(objViewportX, objViewportY, '┃', nil, tcell.StyleDefault)
					case TreeStateSeed:
						screen.SetContent(objViewportX, objViewportY, '.', nil, tcell.StyleDefault)
					case TreeStateAdult:
						screen.SetContent(objViewportX, objViewportY, '█', nil, tcell.StyleDefault)
						screen.SetContent(objViewportX, objViewportY-1, '▄', nil, tcell.StyleDefault)
					}
				}
			}
		}
	}
}

func (game *Game) DrawPlayer(screen tcell.Screen) {
	w, h := screen.Size()
	screen.SetContent(w/2, h/2, CharacterPlayer, nil, tcell.StyleDefault) // Draw the player at the "center" of the view
}

func (game *Game) DrawSquirrel(screen tcell.Screen) {
	w, h := screen.Size()
	screen.SetContent(w/2+game.squirrel.position.x-game.player.position.x, h/2+game.squirrel.position.y-game.player.position.y, CharacterSquirrel, nil, tcell.StyleDefault)

}

func (game *Game) ClearActor(screen tcell.Screen, actorType int) {
	var actor Actor
	switch actorType {
	case ActorPlayer:
		actor = game.player
	case ActorSquirrel:
		actor = game.squirrel
	}
	screen.SetContent(actor.position.x, actor.position.y, ' ', nil, tcell.StyleDefault)
}

func IsBorder(width int, height int, coord Coordinate) (response int, ok bool) {
	/*
		width = 3
		height = 3
		###
		# #
		###

		top border (x = 0..width - 1 (2), y = 0)
		right border (x = width - 1 (2), y = 1..height - 2 (2))
		bottom border (x = 0..width - 1, y = height - 1)
		left border (x = 0, y = 1..height - 2)
	*/
	isTopBorder := (coord.x >= 0 && coord.x <= width-1) && coord.y == 0
	isRightBorder := (coord.x == width-1) && (coord.y >= 1 && coord.y <= height-1)
	isBottomBorder := (coord.x >= 0 && coord.x <= width-1) && coord.y == height
	isLeftBorder := (coord.x == 0) && (coord.y >= 1 && coord.y <= height-1)
	isTopLeftCorner := coord.x == 0 && coord.y == 0
	isTopRightCorner := coord.x == width-1 && coord.y == 0
	isBottomRightCorner := coord.x == width-1 && coord.y == height
	isBottomLeftCorner := coord.x == 0 && coord.y == height

	if isTopBorder {
		if isTopLeftCorner {
			return TopLeftCorner, true
		} else if isTopRightCorner {
			return TopRightCorner, true
		}
		return TopBorder, true
	} else if isRightBorder {
		return RightBorder, true
	} else if isBottomBorder {
		if isBottomRightCorner {
			return BottomRightCorner, true
		} else if isBottomLeftCorner {
			return BottomLeftCorner, true
		}
		return BottomBorder, true
	} else if isLeftBorder {
		return LeftBorder, true
	}

	return -1, false
}
