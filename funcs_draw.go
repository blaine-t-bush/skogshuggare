package main

import "github.com/gdamore/tcell"

var (
	styleTreeWoodNormal = tcell.StyleDefault.Foreground(tcell.NewRGBColor(153, 77, 0))
	styleTreeWoodLight  = tcell.StyleDefault.Foreground(tcell.NewRGBColor(230, 115, 0))
	styleTreeLeaves     = tcell.StyleDefault.Foreground(tcell.NewRGBColor(20, 200, 20))
	styleGrass          = tcell.StyleDefault.Foreground(tcell.NewRGBColor(10, 240, 10))
)

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
			if content, found := game.world.content[Coordinate{x, y}]; found {
				w, h := screen.Size()
				contentViewportX := (w / 2) + (x - game.player.position.x) // Player_viewport_x + Object_real_x - Player_real_x
				contentViewportY := (h / 2) + (y - game.player.position.y) // Player_viewport_y + Object_real_y - Player_real_y
				switch content := content.(type) {
				case Object:
					// Draw object
					if content.collidable {
						screen.SetContent(contentViewportX, contentViewportY, content.char, nil, tcell.StyleDefault)
					} else {
						screen.SetContent(contentViewportX, contentViewportY, content.char, nil, styleGrass)
					}
				case *Tree:
					// Draw tree
					switch content.state {
					case TreeStateStump:
						screen.SetContent(contentViewportX, contentViewportY, '▄', nil, styleTreeWoodNormal)
					case TreeStateTrunk:
						screen.SetContent(contentViewportX, contentViewportY, '█', nil, styleTreeWoodNormal)
					case TreeStateStumpling:
						screen.SetContent(contentViewportX, contentViewportY, '╻', nil, styleTreeWoodLight)
					case TreeStateSapling:
						screen.SetContent(contentViewportX, contentViewportY, '┃', nil, styleTreeWoodLight)
					case TreeStateSeed:
						screen.SetContent(contentViewportX, contentViewportY, '.', nil, styleTreeWoodLight)
					case TreeStateAdult:
						screen.SetContent(contentViewportX, contentViewportY, '█', nil, styleTreeWoodNormal)
						screen.SetContent(contentViewportX-1, contentViewportY-1, '▓', nil, styleTreeLeaves)
						screen.SetContent(contentViewportX, contentViewportY-1, '▓', nil, styleTreeLeaves)
						screen.SetContent(contentViewportX+1, contentViewportY-1, '▓', nil, styleTreeLeaves)
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
