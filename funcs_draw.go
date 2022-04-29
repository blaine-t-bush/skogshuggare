package main

import (
	"strconv"

	"github.com/gdamore/tcell"
)

func (game *Game) Draw() {
	game.screen.Clear()
	game.DrawViewport()
	game.DrawMenu()
	game.screen.Show()
}

// Only draw things within the player view range
// Draw the player last, run checks in DrawPlayer function to check if player should be drawn or not.
func (game *Game) DrawViewport() {
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
	w, h := game.screen.Size()
	playerViewportCoord := Coordinate{w / 2, h / 2}

	// Get squirrel-blocked tiles.
	var squirrelViewportCoord Coordinate
	var squirrelViewportCoords []Coordinate
	for _, squirrel := range game.squirrels {
		squirrelViewportCoord = Translate(playerViewportCoord, squirrel.position.x-game.player.position.x, squirrel.position.y-game.player.position.y)
		squirrelViewportCoords = append(squirrelViewportCoords, squirrelViewportCoord)
	}

	// Draw content.
	aboveViewportCoords := []Coordinate{}
	xRadiusMin, xRadiusMax, yRadiusMin, yRadiusMax := game.GetDrawRanges()
	for x := xRadiusMin; x <= xRadiusMax; x++ {
		for y := yRadiusMin; y <= yRadiusMax; y++ {
			coord := Coordinate{x, y}

			// Get the viewport coordinates
			contentViewportCoord := Translate(playerViewportCoord, x-game.player.position.x, y-game.player.position.y)

			if border, isBorder := game.world.borders[coord]; isBorder {
				switch border {
				case TopBorder, BottomBorder:
					game.screen.SetContent(contentViewportCoord.x, contentViewportCoord.y, tcell.RuneHLine, nil, tcell.StyleDefault)
				case RightBorder, LeftBorder:
					game.screen.SetContent(contentViewportCoord.x, contentViewportCoord.y, tcell.RuneVLine, nil, tcell.StyleDefault)
				case TopLeftCorner:
					game.screen.SetContent(contentViewportCoord.x, contentViewportCoord.y, tcell.RuneULCorner, nil, tcell.StyleDefault)
				case TopRightCorner:
					game.screen.SetContent(contentViewportCoord.x, contentViewportCoord.y, tcell.RuneURCorner, nil, tcell.StyleDefault)
				case BottomRightCorner:
					game.screen.SetContent(contentViewportCoord.x, contentViewportCoord.y, tcell.RuneLRCorner, nil, tcell.StyleDefault)
				case BottomLeftCorner:
					game.screen.SetContent(contentViewportCoord.x, contentViewportCoord.y, tcell.RuneLLCorner, nil, tcell.StyleDefault)
				}
				continue
			}

			if content, found := game.world.content[Coordinate{x, y}]; found {
				switch content := content.(type) {
				case *StaticObject:
					// Draw object
					aboveViewportCoords = game.DrawContent(content.key, contentViewportCoord, aboveViewportCoords)
				case *AnimatedObject:
					// Draw object
					aboveViewportCoords = game.DrawContent(GetAnimationState(content.key, content.animationStage), contentViewportCoord, aboveViewportCoords)
				case *Fire:
					aboveViewportCoords = game.DrawContent(GetAnimationState(KeyFire, content.animationStage), contentViewportCoord, aboveViewportCoords)
				case *Tree:
					// Draw tree
					switch content.state {
					case TreeStateStump:
						aboveViewportCoords = game.DrawContent(KeyTreeStump, contentViewportCoord, aboveViewportCoords)
					case TreeStateTrunk:
						aboveViewportCoords = game.DrawContent(KeyTreeTrunk, contentViewportCoord, aboveViewportCoords)
					case TreeStateStumpling:
						aboveViewportCoords = game.DrawContent(KeyTreeStumpling, contentViewportCoord, aboveViewportCoords)
					case TreeStateSapling:
						aboveViewportCoords = game.DrawContent(KeyTreeSapling, contentViewportCoord, aboveViewportCoords)
					case TreeStateSeed:
						aboveViewportCoords = game.DrawContent(KeyTreeSeed, contentViewportCoord, aboveViewportCoords)
					case TreeStateAdult:
						aboveViewportCoords = game.DrawContent(KeyTreeTrunk, contentViewportCoord, aboveViewportCoords)
						aboveViewportCoords = game.DrawContent(KeyTreeLeaves, Translate(contentViewportCoord, -1, -1), aboveViewportCoords)
						aboveViewportCoords = game.DrawContent(KeyTreeLeaves, Translate(contentViewportCoord, 0, -1), aboveViewportCoords)
						aboveViewportCoords = game.DrawContent(KeyTreeLeaves, Translate(contentViewportCoord, 1, -1), aboveViewportCoords)
					}
				}
			}

			if decoration, found := game.world.decorations[Coordinate{x, y}]; found {
				switch decoration := decoration.(type) {
				case *DecorationObject:
					aboveViewportCoords = game.DrawContent(GetAnimationState(decoration.key, decoration.animationStage), contentViewportCoord, aboveViewportCoords)
				}
			}
		}
	}

	game.DrawActor(KeyPlayer, playerViewportCoord, aboveViewportCoords)
	for _, squirrelViewportCoord := range squirrelViewportCoords {
		game.DrawActor(KeySquirrel, squirrelViewportCoord, aboveViewportCoords) // FIXME only draw inside viewport
	}
}

// Draws content for the given key at the given coord, and updates the list of coordinates to render above actors if necessary.
func (game *Game) DrawContent(key int, coord Coordinate, aboveCoords []Coordinate) []Coordinate {
	symbol := symbols[key]
	game.screen.SetContent(coord.x, coord.y, symbol.char, nil, symbol.style)
	if symbol.aboveActor {
		aboveCoords = append(aboveCoords, coord)
	}

	return aboveCoords
}

// Draws actor, but only if they are not on a tile with priority content (i.e. aboveActor is true).
// Also gives actor the same background as the content they're on, if it exists.
func (game *Game) DrawActor(key int, coord Coordinate, aboveCoords []Coordinate) bool {
	symbol := symbols[key]
	hasAbove := false
	for _, aboveCoord := range aboveCoords {
		if aboveCoord == coord {
			hasAbove = true
			break
		}
	}

	var bg tcell.Color
	hasBelow := false
	w, h := game.screen.Size()
	playerViewportCoord := Coordinate{w / 2, h / 2}
	worldCoord := Translate(coord, -playerViewportCoord.x+game.player.position.x, -playerViewportCoord.y+game.player.position.y)
	if content, exists := game.world.content[worldCoord]; exists {
		switch content := content.(type) {
		case *StaticObject:
			_, bg, _ = symbols[content.key].style.Decompose()
			hasBelow = true
		case *AnimatedObject:
			_, bg, _ = symbols[content.key].style.Decompose()
			hasBelow = true
		case *Fire:
			_, bg, _ = symbols[KeyFire].style.Decompose()
			hasBelow = true
		}
	}

	if !hasAbove {
		if hasBelow {
			game.screen.SetContent(coord.x, coord.y, symbol.char, nil, symbol.style.Background(bg))
		} else {
			game.screen.SetContent(coord.x, coord.y, symbol.char, nil, symbol.style)
		}
		return true
	}

	return false
}

func (game *Game) DrawMenu() {
	game.DrawMenuBorder()
	// Draw score: 0
	//      12345678
	scoreString := "Score: " + strconv.Itoa(game.player.score)
	scoreIdx := 0
	for i := 1; i < len(scoreString)+1; i++ {
		game.screen.SetContent(i, 1, rune(scoreString[scoreIdx]), nil, tcell.StyleDefault)
		scoreIdx++
	}

	hitPointsString := "HP: " + strconv.Itoa(game.player.hitPointsCurrent)
	hitPointsIdx := 0
	for i := 1; i < len(hitPointsString)+1; i++ {
		game.screen.SetContent(i, 2, rune(hitPointsString[hitPointsIdx]), nil, tcell.StyleDefault)
		hitPointsIdx++
	}

	game.PrintToMenu()
}

func (game *Game) DrawMenuBorder() {
	for c := 1; c < game.menu.width; c++ { // Draw top and bottom borders
		game.screen.SetContent(c, 0, tcell.RuneHLine, nil, tcell.StyleDefault)
		game.screen.SetContent(c, game.menu.height, tcell.RuneHLine, nil, tcell.StyleDefault)
	}

	for r := 1; r <= game.menu.height-1; r++ { // Add left and right borders
		game.screen.SetContent(0, r, tcell.RuneVLine, nil, tcell.StyleDefault)
		game.screen.SetContent(game.menu.width, r, tcell.RuneVLine, nil, tcell.StyleDefault)
	}

	// Add corners
	game.screen.SetContent(0, 0, tcell.RuneULCorner, nil, tcell.StyleDefault)
	game.screen.SetContent(game.menu.width, 0, tcell.RuneURCorner, nil, tcell.StyleDefault)
	game.screen.SetContent(0, game.menu.height, tcell.RuneLLCorner, nil, tcell.StyleDefault)
	game.screen.SetContent(game.menu.width, game.menu.height, tcell.RuneLRCorner, nil, tcell.StyleDefault)
}

func (game *Game) PrintToMenu() {
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
			game.screen.SetContent(currX, currY, r, nil, tcell.StyleDefault)
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
	xRadiusMax = game.world.width - 1
	yRadiusMin = 0
	yRadiusMax = game.world.height - 1

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
