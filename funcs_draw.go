package main

import (
	"strconv"

	"github.com/gdamore/tcell"
)

func (game *Game) Draw(screen tcell.Screen) {
	screen.Clear()
	game.DrawViewport(screen)
	game.DrawMenu(screen)
	screen.Show()
}

// Only draw things within the player view range
// Draw the player last, run checks in DrawPlayer function to check if player should be drawn or not.
func (game *Game) DrawViewport(screen tcell.Screen) {
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

	// Draw player.
	w, h := screen.Size()
	playerViewportCoord := Coordinate{w / 2, h / 2}
	game.DrawContent(screen, KeyPlayer, playerViewportCoord, []Coordinate{})

	// Draw squirrels.
	var squirrelViewportCoord Coordinate
	var squirrelViewportCoords []Coordinate
	for _, squirrel := range game.squirrels {
		squirrelViewportCoord = Coordinate{playerViewportCoord.x + squirrel.position.x - game.player.position.x, playerViewportCoord.y + squirrel.position.y - game.player.position.y}
		game.DrawContent(screen, KeySquirrel, squirrelViewportCoord, []Coordinate{playerViewportCoord}) // FIXME only draw inside viewport
		squirrelViewportCoords = append(squirrelViewportCoords, squirrelViewportCoord)
	}

	// Draw content.
	actorViewportCoords := append(squirrelViewportCoords, playerViewportCoord)
	xRadiusMin, xRadiusMax, yRadiusMin, yRadiusMax := game.GetDrawRanges()
	for x := xRadiusMin; x <= xRadiusMax; x++ {
		for y := yRadiusMin; y <= yRadiusMax; y++ {
			coord := Coordinate{x, y}

			// Get the viewport coordinates
			contentViewportX := playerViewportCoord.x + (x - game.player.position.x) // Player_viewport_x + Object_real_x - Player_real_x
			contentViewportY := playerViewportCoord.y + (y - game.player.position.y) // Player_viewport_y + Object_real_y - Player_real_y
			contentViewportCoord := Coordinate{contentViewportX, contentViewportY}

			if border, isBorder := game.world.borders[coord]; isBorder {
				switch border {
				case TopBorder, BottomBorder:
					screen.SetContent(contentViewportX, contentViewportY, tcell.RuneHLine, nil, tcell.StyleDefault)
				case RightBorder, LeftBorder:
					screen.SetContent(contentViewportX, contentViewportY, tcell.RuneVLine, nil, tcell.StyleDefault)
				case TopLeftCorner:
					screen.SetContent(contentViewportX, contentViewportY, tcell.RuneULCorner, nil, tcell.StyleDefault)
				case TopRightCorner:
					screen.SetContent(contentViewportX, contentViewportY, tcell.RuneURCorner, nil, tcell.StyleDefault)
				case BottomRightCorner:
					screen.SetContent(contentViewportX, contentViewportY, tcell.RuneLRCorner, nil, tcell.StyleDefault)
				case BottomLeftCorner:
					screen.SetContent(contentViewportX, contentViewportY, tcell.RuneLLCorner, nil, tcell.StyleDefault)
				}
				continue
			}

			if content, found := game.world.content[Coordinate{x, y}]; found {
				switch content := content.(type) {
				case Object:
					// Draw object
					game.DrawContent(screen, content.key, contentViewportCoord, actorViewportCoords)
				case *Tree:
					// Draw tree
					switch content.state {
					case TreeStateStump:
						game.DrawContent(screen, KeyTreeStump, contentViewportCoord, actorViewportCoords)
					case TreeStateTrunk:
						game.DrawContent(screen, KeyTreeTrunk, contentViewportCoord, actorViewportCoords)
					case TreeStateStumpling:
						game.DrawContent(screen, KeyTreeStumpling, contentViewportCoord, actorViewportCoords)
					case TreeStateSapling:
						game.DrawContent(screen, KeyTreeSapling, contentViewportCoord, actorViewportCoords)
					case TreeStateSeed:
						game.DrawContent(screen, KeyTreeSeed, contentViewportCoord, actorViewportCoords)
					case TreeStateAdult:
						game.DrawContent(screen, KeyTreeTrunk, contentViewportCoord, actorViewportCoords)
						game.DrawContent(screen, KeyTreeLeaves, Coordinate{contentViewportCoord.x - 1, contentViewportCoord.y - 1}, actorViewportCoords)
						game.DrawContent(screen, KeyTreeLeaves, Coordinate{contentViewportCoord.x, contentViewportCoord.y - 1}, actorViewportCoords)
						game.DrawContent(screen, KeyTreeLeaves, Coordinate{contentViewportCoord.x + 1, contentViewportCoord.y - 1}, actorViewportCoords)
					}
				}
			}
		}
	}
}

// Draws content for the given key at the given coord, but only if that coord is not in priorityCoords
func (game *Game) DrawContent(screen tcell.Screen, key int, coord Coordinate, priorityCoords []Coordinate) {
	symbol := symbols[key]
	draw := true
	for _, priorityCoord := range priorityCoords {
		if coord == priorityCoord && !symbol.aboveActor {
			draw = false
		}
	}

	if draw {
		screen.SetContent(coord.x, coord.y, symbol.char, nil, symbol.style)
	}
}

func (game *Game) DrawMenu(screen tcell.Screen) {
	game.DrawMenuBorder(screen)
	// Draw score: 0
	//      12345678
	scoreString := "Score: " + strconv.Itoa(game.player.score)
	scoreIdx := 0
	for i := 1; i < len(scoreString)+1; i++ {
		screen.SetContent(i, 1, rune(scoreString[scoreIdx]), nil, tcell.StyleDefault)
		scoreIdx++
	}
	game.PrintToMenu(screen)
}

func (game *Game) DrawMenuBorder(screen tcell.Screen) {
	for c := 1; c < game.menu.width; c++ { // Draw top and bottom borders
		screen.SetContent(c, 0, tcell.RuneHLine, nil, tcell.StyleDefault)
		screen.SetContent(c, game.menu.height, tcell.RuneHLine, nil, tcell.StyleDefault)
	}

	for r := 1; r <= game.menu.height-1; r++ { // Add left and right borders
		screen.SetContent(0, r, tcell.RuneVLine, nil, tcell.StyleDefault)
		screen.SetContent(game.menu.width, r, tcell.RuneVLine, nil, tcell.StyleDefault)
	}

	// Add corners
	screen.SetContent(0, 0, tcell.RuneULCorner, nil, tcell.StyleDefault)
	screen.SetContent(game.menu.width, 0, tcell.RuneURCorner, nil, tcell.StyleDefault)
	screen.SetContent(0, game.menu.height, tcell.RuneLLCorner, nil, tcell.StyleDefault)
	screen.SetContent(game.menu.width, game.menu.height, tcell.RuneLRCorner, nil, tcell.StyleDefault)
}

func (game *Game) PrintToMenu(screen tcell.Screen) {
	maxLen := game.menu.width
	maxHeight := game.menu.height

	currX := 1
	currY := 2
	for _, message := range game.menu.messages {
		for c := 0; c < len(message); c++ {
			r := rune(message[c])
			if c%maxLen == 0 {
				currX = 1
				currY++
			}
			if currY >= maxHeight {
				break
			}
			screen.SetContent(currX, currY, r, nil, tcell.StyleDefault)
			currX++
		}
	}

}

func (game *Game) AppendToMenuMessages(text string) {
	if len(game.menu.messages) <= 2 {
		game.menu.messages = append(game.menu.messages, text)
	} else {
		game.menu.messages = append(game.menu.messages[:0], game.menu.messages[1:]...)
	}
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
	isRightBorder := coord.x == width-1 && (coord.y >= 0 && coord.y <= height-1)
	isBottomBorder := (coord.x >= 0 && coord.x <= width-1) && coord.y == height-1
	isLeftBorder := coord.x == 0 && (coord.y >= 0 && coord.y <= height-1)
	isTopLeftCorner := coord.x == 0 && coord.y == 0
	isTopRightCorner := coord.x == width-1 && coord.y == 0
	isBottomRightCorner := coord.x == width-1 && coord.y == height-1
	isBottomLeftCorner := coord.x == 0 && coord.y == height-1

	if isTopLeftCorner {
		return TopLeftCorner, true
	} else if isTopRightCorner {
		return TopRightCorner, true
	} else if isBottomRightCorner {
		return BottomRightCorner, true
	} else if isBottomLeftCorner {
		return BottomLeftCorner, true
	} else if isTopBorder {
		return TopBorder, true
	} else if isRightBorder {
		return RightBorder, true
	} else if isBottomBorder {
		return BottomBorder, true
	} else if isLeftBorder {
		return LeftBorder, true
	}

	return -1, false
}

func (game *Game) GetDrawRanges() (xRadiusMin int, xRadiusMax int, yRadiusMin int, yRadiusMax int) {
	xRadiusMin = 0
	xRadiusMax = game.world.width
	yRadiusMin = 0
	yRadiusMax = game.world.height

	if game.player.position.x-game.player.visionRadius > 0 {
		xRadiusMin = game.player.position.x - game.player.visionRadius
	}

	if game.player.position.x+game.player.visionRadius < game.world.width {
		xRadiusMax = game.player.position.x + game.player.visionRadius
	}

	if game.player.position.y-game.player.visionRadius > 0 {
		yRadiusMin = game.player.position.y - game.player.visionRadius
	}

	if game.player.position.y+game.player.visionRadius < game.world.height {
		yRadiusMax = game.player.position.y + game.player.visionRadius
	}

	return xRadiusMin, xRadiusMax, yRadiusMin, yRadiusMax
}
