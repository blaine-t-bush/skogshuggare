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
			obj, found := game.world.content[Coordinate{x, y}]
			if found {
				w, h := screen.Size()
				objViewportX := (w / 2) + (x - game.player.position.x) // Player_viewport_x + Object_real_x - Player_real_x
				objViewportY := (h / 2) + (y - game.player.position.y) // Player_viewport_y + Object_real_y - Player_real_y
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
