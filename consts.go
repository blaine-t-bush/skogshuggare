package main

const (
	// Game parameters
	TickRate = 30 // Milliseconds between ticks
	// Actors
	ActorPlayer   = 1
	ActorSquirrel = 2
	// Characters
	CharacterPlayer   = '@'
	CharacterSquirrel = '~'
	// Directions
	DirUp     = 0
	DirRight  = 1
	DirDown   = 2
	DirLeft   = 3
	DirOmni   = 4
	DirRandom = 5
	DirNone   = 6
	// Living tree states
	TreeStateSeed    = 0
	TreeStateSapling = 1
	TreeStateAdult   = 2
	// Harvested tree states
	TreeStateRemoved   = 10
	TreeStateStump     = 11
	TreeStateTrunk     = 12
	TreeStateStumpling = 13
	// Growth chances (per game tick)
	GrowthChanceSeed    = 0.010 // Seed to sapling
	GrowthChanceSapling = 0.005 // Sapling to adult
	SeedCreationChance  = 0.005 // Seed spawning
	SeedCreationMax     = 3     // Maximum number of seeds to create per tick
	// Border states
	TopBorder         = 100
	RightBorder       = 101
	BottomBorder      = 102
	LeftBorder        = 103
	TopLeftCorner     = 104
	TopRightCorner    = 105
	BottomRightCorner = 106
	BottomLeftCorner  = 107
)

// NOTE
// Interesting Unicode characters (e.g. arrows) start at 2190.
