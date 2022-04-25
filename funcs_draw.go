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

	w, h := screen.Size()
	playerViewportCoord := Coordinate{w / 2, h / 2}
	game.DrawPlayer(screen)

	for x := game.player.position.x - game.player.visionRadius; x <= game.player.position.x+game.player.visionRadius; x++ {
		for y := game.player.position.y - game.player.visionRadius; y <= game.player.position.y+game.player.visionRadius; y++ {
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
					game.DrawContent(screen, content.key, contentViewportCoord, playerViewportCoord)
				case *Tree:
					// Draw tree
					switch content.state {
					case TreeStateStump:
						game.DrawContent(screen, KeyTreeStump, contentViewportCoord, playerViewportCoord)
					case TreeStateTrunk:
						game.DrawContent(screen, KeyTreeTrunk, contentViewportCoord, playerViewportCoord)
					case TreeStateStumpling:
						game.DrawContent(screen, KeyTreeStumpling, contentViewportCoord, playerViewportCoord)
					case TreeStateSapling:
						game.DrawContent(screen, KeyTreeSapling, contentViewportCoord, playerViewportCoord)
					case TreeStateSeed:
						game.DrawContent(screen, KeyTreeSeed, contentViewportCoord, playerViewportCoord)
					case TreeStateAdult:
						game.DrawContent(screen, KeyTreeTrunk, contentViewportCoord, playerViewportCoord)
						game.DrawContent(screen, KeyTreeLeaves, Coordinate{contentViewportCoord.x - 1, contentViewportCoord.y - 1}, playerViewportCoord)
						game.DrawContent(screen, KeyTreeLeaves, Coordinate{contentViewportCoord.x, contentViewportCoord.y - 1}, playerViewportCoord)
						game.DrawContent(screen, KeyTreeLeaves, Coordinate{contentViewportCoord.x + 1, contentViewportCoord.y - 1}, playerViewportCoord)
					}
				}
			}

			if coord == game.squirrel.position {
				game.DrawContent(screen, KeySquirrel, contentViewportCoord, playerViewportCoord)
			}
		}
	}
}

func (game *Game) DrawPlayer(screen tcell.Screen) {
	w, h := screen.Size()
	playerX, playerY := w/2, h/2                                                                // Save player position to these variables. Previous implementation caused variables to be changed multiple times.
	screen.SetContent(playerX, playerY, symbols[KeyPlayer].char, nil, symbols[KeyPlayer].style) // Draw the player at the "center" of the view
}

func (game *Game) DrawContent(screen tcell.Screen, key int, coord Coordinate, playerCoord Coordinate) {
	symbol := symbols[key]
	if coord != playerCoord || symbol.aboveActor {
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
	if len(game.menu.messages) <= 3 {
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
