package main

import "github.com/gdamore/tcell"

const (
	// Game parameters
	TickRate = 30 // Milliseconds between ticks
	// Actors
	ActorPlayer   = 1
	ActorSquirrel = 2
	// Runes
	RunePlayer        = '@'
	RuneSquirrel      = 's'
	RuneWall          = '#'
	RuneTreeSeed      = '.'
	RuneTreeSapling   = '┃'
	RuneTreeTrunk     = '█'
	RuneTreeLeaves    = '▓'
	RuneTreeStump     = '▄'
	RuneTreeStumpling = '╻'
	RuneGrassLight    = '\''
	RuneGrassHeavy    = '"'
	RuneWater         = '~'
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

var (
	symbolColors = map[rune]tcell.Color{
		RunePlayer:        tcell.ColorIndianRed,
		RuneSquirrel:      tcell.ColorRosyBrown,
		RuneWall:          tcell.ColorWhite,
		RuneTreeSeed:      tcell.ColorKhaki,
		RuneTreeSapling:   tcell.ColorSaddleBrown,
		RuneTreeTrunk:     tcell.ColorSaddleBrown,
		RuneTreeLeaves:    tcell.ColorForestGreen,
		RuneTreeStump:     tcell.ColorSaddleBrown,
		RuneTreeStumpling: tcell.ColorSaddleBrown,
		RuneGrassLight:    tcell.ColorGreenYellow,
		RuneGrassHeavy:    tcell.ColorGreenYellow,
		RuneWater:         tcell.ColorCornflowerBlue,
	}
	symbolStyles = map[rune]tcell.Style{
		RunePlayer:        tcell.StyleDefault.Foreground(symbolColors[RunePlayer]),
		RuneSquirrel:      tcell.StyleDefault.Foreground(symbolColors[RuneSquirrel]),
		RuneWall:          tcell.StyleDefault.Foreground(symbolColors[RuneWall]),
		RuneTreeSeed:      tcell.StyleDefault.Foreground(symbolColors[RuneTreeSeed]),
		RuneTreeSapling:   tcell.StyleDefault.Foreground(symbolColors[RuneTreeSapling]),
		RuneTreeTrunk:     tcell.StyleDefault.Foreground(symbolColors[RuneTreeTrunk]),
		RuneTreeLeaves:    tcell.StyleDefault.Foreground(symbolColors[RuneTreeLeaves]),
		RuneTreeStump:     tcell.StyleDefault.Foreground(symbolColors[RuneTreeStump]),
		RuneTreeStumpling: tcell.StyleDefault.Foreground(symbolColors[RuneTreeStumpling]),
		RuneGrassLight:    tcell.StyleDefault.Foreground(symbolColors[RuneGrassLight]),
		RuneGrassHeavy:    tcell.StyleDefault.Foreground(symbolColors[RuneGrassHeavy]),
		RuneWater:         tcell.StyleDefault.Foreground(symbolColors[RuneWater]),
	}
)

// NOTE
// Interesting Unicode characters (e.g. arrows) start at 2190.
