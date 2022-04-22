package main

type Coordinate struct {
	x int
	y int
}

type Actor struct {
	position     Coordinate
	visionRadius int
}

type Border struct {
	x1 int // Left border boundary x-coordinate
	x2 int // Right border boundary x-coordinate
	y1 int // Top border boundary y-coordinate
	y2 int // Bottom border boundary y-coordinate
	t  int // Border thickness in characters
}

type Tree struct {
	position Coordinate
	state    int // See constants
}

type Object struct {
	char       rune
	collidable bool
}

type World struct {
	width   int
	height  int
	borders map[Coordinate]int // Store the borders in a lookup-table instead of running checks every single loop
	content map[Coordinate]interface{}
}

type Game struct {
	player   Actor
	squirrel Actor
	border   Border
	world    World
	exit     bool
}
