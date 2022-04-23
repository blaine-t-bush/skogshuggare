package main

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
	borders map[Coordinate]int // Store the borders in a lookup-table instead of running checks every single loop
	content map[Coordinate]interface{}
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