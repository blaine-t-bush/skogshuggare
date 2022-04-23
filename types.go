package main

type Coordinate struct {
	x int
	y int
}

type Actor struct {
	position     Coordinate
	destination  Coordinate
	visionRadius int
}

type Tree struct {
	position Coordinate
	state    int // See constants
}

type Object struct {
	char       rune
	collidable bool  // Are actors blocked
	plantable  bool  // Can seeds be planted here
	r          int32 // Rgb color value
	g          int32 // rGb color value
	b          int32 // rgB color value
}

type World struct {
	width   int
	height  int
	content map[Coordinate]interface{}
}

type Game struct {
	player   Actor
	squirrel Actor
	world    World
	exit     bool
}
