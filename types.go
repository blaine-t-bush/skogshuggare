package main

import "github.com/gdamore/tcell"

type Coordinate struct {
	x int
	y int
}

type Actor struct {
	position     Coordinate
	destination  Coordinate
	visionRadius int
	score        int
}

type Tree struct {
	position Coordinate
	state    int // See constants
}

type Object struct {
	key        int
	collidable bool // Are actors blocked
	plantable  bool // Can seeds be planted here
}

type World struct {
	width   int
	height  int
	borders map[Coordinate]int // Store the borders in a lookup-table instead of running checks every single loop
	content map[Coordinate]any
}

type Game struct {
	player   Actor
	squirrel Actor
	world    World
	menu     Menu
	exit     bool
}

type Menu struct {
	width    int
	height   int
	position Coordinate
	messages []string
}

type Symbol struct {
	char  rune
	style tcell.Style
}

type GrowthInfo struct {
	newState int
	chance   float64
}

type TitleMenu struct {
	cursorState    int
	pageState      int
	titleMenuPages map[int]*TitleMenuPage
	exit           bool
}

type TitleMenuItem struct {
	order int
	text  string
	value interface{}
}

type TitleMenuPage struct {
	name           int
	content        []string
	animationState int
	cursorState    int
	titleMenuItems map[int]TitleMenuItem
}
