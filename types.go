package main

import "github.com/gdamore/tcell"

type Coordinate struct {
	x int
	y int
}

type Actor struct {
	position     Coordinate
	destination  Coordinate
	path         map[int]Coordinate
	visionRadius int
	score        int
}

type Tree struct {
	position Coordinate
	state    int // See constants
}

type Object struct {
	key        int
	category   int
	collidable bool // Are actors blocked
	plantable  bool // Can seeds be planted here
}

type World struct {
	width   int
	height  int
	borders map[Coordinate]int // Store the borders in a lookup-table instead of running checks every single loop
	content map[Coordinate]interface{}
}

type Game struct {
	player    Actor
	squirrels map[int]*Actor
	world     World
	menu      Menu
	exit      bool
}

type Menu struct {
	width    int
	height   int
	position Coordinate
	messages []string
}

type Symbol struct {
	char       rune
	aboveActor bool
	style      tcell.Style
}

type GrowthInfo struct {
	newState int
	chance   float64
}
