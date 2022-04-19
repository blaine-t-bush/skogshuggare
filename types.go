package main

import "github.com/gdamore/tcell"

type Player struct {
	x int // Player x-coordinate
	y int // Player y-coordinate
}

type Border struct {
	x1 int // Left border boundary x-coordinate
	x2 int // Right border boundary x-coordinate
	y1 int // Top border boundary y-coordinate
	y2 int // Bottom border boundary y-coordinate
	t  int // Border thickness in characters
}

type Tree struct {
	x int // Trunk left corner x-coordinate
	y int // Trunk left corner y-coordinate
}

type Game struct {
	player Player
	border Border
	trees  []Tree
}

func (game *Game) DrawBorder(screen tcell.Screen) {
	for c := game.border.x1; c <= game.border.x2; c++ { // Add top and bottom borders
		screen.SetContent(c, game.border.y1, '#', nil, tcell.StyleDefault)
		screen.SetContent(c, game.border.y2, '#', nil, tcell.StyleDefault)
	}

	for r := game.border.y1 + 1; r <= game.border.y2-1; r++ { // Add left and right borders
		screen.SetContent(game.border.x1, r, '#', nil, tcell.StyleDefault)
		screen.SetContent(game.border.x2, r, '#', nil, tcell.StyleDefault)
	}
}

func (game *Game) DrawTrees(screen tcell.Screen) {
	//  /\
	// /__\
	//  ||
	for _, tree := range game.trees {
		screen.SetContent(tree.x, tree.y, '|', nil, tcell.StyleDefault)
		screen.SetContent(tree.x+1, tree.y, '|', nil, tcell.StyleDefault)
		screen.SetContent(tree.x-1, tree.y-1, '/', nil, tcell.StyleDefault)
		screen.SetContent(tree.x, tree.y-1, '_', nil, tcell.StyleDefault)
		screen.SetContent(tree.x+1, tree.y-1, '_', nil, tcell.StyleDefault)
		screen.SetContent(tree.x+2, tree.y-1, '\\', nil, tcell.StyleDefault)
		screen.SetContent(tree.x, tree.y-2, '/', nil, tcell.StyleDefault)
		screen.SetContent(tree.x+1, tree.y-2, '\\', nil, tcell.StyleDefault)
	}
}

func (game *Game) ClearPlayer(screen tcell.Screen) {
	screen.SetContent(game.player.x, game.player.y, ' ', nil, tcell.StyleDefault)
	// screen.Show()
}

func (game *Game) DrawPlayer(screen tcell.Screen) {
	screen.SetContent(game.player.x, game.player.y, '@', nil, tcell.StyleDefault)
	// screen.Show()
}

func (game *Game) Draw(screen tcell.Screen) {
	game.DrawBorder(screen)
	game.DrawPlayer(screen)
	game.DrawTrees(screen)
	screen.Show()
}

func (game *Game) MovePlayer(screen tcell.Screen, len int, dir int) {
	game.ClearPlayer(screen)

	// Determine (potential) new location.
	pMoved := game.player
	if len != 0 {
		switch dir {
		case 0: // up
			pMoved.y = game.player.y - len
		case 1: // right
			pMoved.x = game.player.x + len
		case 2: // down
			pMoved.y = game.player.y + len
		case 3: // left
			pMoved.x = game.player.x - len
		}
	}

	// TODO
	// Prevent update if player would collide with tree trunks, but not with
	// tree canopies.
	// Trunks are located at tree.x, tree.y and tree.x+1, tree.y
	for _, tree := range game.trees {
		if pMoved.y == tree.y && (pMoved.x == tree.x || pMoved.x == tree.x+1) {
			pMoved.x = game.player.x
			pMoved.y = game.player.y
		}
	}

	// Prevent update if new location is past left or right boundaries.
	if pMoved.x <= game.border.x1 {
		pMoved.x = game.border.x1 + game.border.t
	} else if pMoved.x >= game.border.x2 {
		pMoved.x = game.border.x2 - game.border.t
	}

	// Prevent update if new location is past top or bottom boundaries.
	if pMoved.y <= game.border.y1 {
		pMoved.y = game.border.y1 + game.border.t
	} else if pMoved.y >= game.border.y2 {
		pMoved.y = game.border.y2 - game.border.t
	}

	game.player = pMoved
	game.Draw(screen)
}
